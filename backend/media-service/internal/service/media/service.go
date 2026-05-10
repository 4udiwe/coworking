package media_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
	object_repository "github.com/4udiwe/coworking/backend/media-service/internal/repository/object"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
)

type MediaService struct {
	repo      MediaRepository
	storage   ObjectStorage
	processor ImageProcessor

	urlTTL time.Duration
}

func New(repo MediaRepository, storage ObjectStorage, processor ImageProcessor) *MediaService {
	return &MediaService{
		repo:      repo,
		storage:   storage,
		processor: processor,
		urlTTL:    15 * time.Minute,
	}
}

func (s *MediaService) Upload(ctx context.Context, input entity.UploadInput) (entity.UploadResult, error) {

	// 3. Создаём media
	media := entity.Media{
		OriginalName: input.FileName,
		MimeType:     input.ContentType,
		Status:       entity.StatusPending,
	}

	id, err := s.repo.Create(ctx, media)
	if err != nil {
		return entity.UploadResult{}, err
	}

	media.ID = id

	// 4. Upload original
	originalKey := object_repository.BuildKey(media.ID.Hex(), string(entity.SizeOriginal))

	if err := s.storage.Upload(ctx, originalKey, input.Data, "image/webp"); err != nil {
		return entity.UploadResult{}, err
	}

	// 4.1 sync resize для thumbnail
	thumbData, _, _, err := s.processor.ResizeToWidth(ctx, input.Data, entity.SizeWidths[entity.SizeThumbnail])
	if err != nil {
		return entity.UploadResult{}, err
	}
	thumbKey := object_repository.BuildKey(media.ID.Hex(), string(entity.SizeThumbnail))

	// 5. Upload thumbnail
	if err := s.storage.Upload(ctx, thumbKey, thumbData, "image/webp"); err != nil {
		return entity.UploadResult{}, err
	}

	thumbVariant := entity.ImageVariant{
		Size:       entity.SizeThumbnail,
		Width:      entity.SizeWidths[entity.SizeThumbnail],
		StorageKey: thumbKey,
	}

	// 6. статус + variant
	if err := s.repo.UpdateStatusAndVariants(
		ctx,
		id,
		entity.StatusProcessing, // Устанавливаем статус "processing", после завершения async-обработки обновим на "ready"
		[]entity.ImageVariant{thumbVariant},
	); err != nil {
		return entity.UploadResult{}, err
	}

	// 7. async resize
	s.processAsyncSizes(media)

	return entity.UploadResult{
		ID:     id.Hex(),
		Status: string(entity.StatusProcessing),
		URLs: map[string]string{
			"thumbnail": thumbKey,
		},
	}, nil
}

func (s *MediaService) processAsyncSizes(media entity.Media) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		g, ctx := errgroup.WithContext(ctx)

		log := logrus.WithField("media_id", media.ID.Hex())

		// скачиваем original
		originalKey := object_repository.BuildKey(media.ID.Hex(), string(entity.SizeOriginal))

		originalData, err := s.storage.Get(ctx, originalKey)
		if err != nil {
			log.WithError(err).Error("failed to download original")
			return
		}

		var (
			mu       sync.Mutex
			variants []entity.ImageVariant
		)

		// 1. запускаем горутину на каждый size
		for _, rsize := range entity.AsyncSizes {

			size := rsize

			g.Go(func() error {

				// resize
				resized, _, _, err := s.processor.ResizeToWidth(ctx, originalData, entity.SizeWidths[size])
				if err != nil {
					return err
				}

				key := object_repository.BuildKey(media.ID.Hex(), string(size))

				// upload
				err = s.storage.Upload(ctx, key, resized, "image/webp")
				if err != nil {
					return err
				}

				v := entity.ImageVariant{
					Size:       size,
					Width:      entity.SizeWidths[size],
					StorageKey: key,
					Bytes:      int64(len(resized)),
				}

				mu.Lock()
				variants = append(variants, v)
				mu.Unlock()

				return nil
			})
		}

		// 2. ждём все горутины
		err = g.Wait()

		// если хотя бы одна упала — НЕ фиксируем ready
		if err != nil {
			log.WithError(err).Error("processAsync failed, leaving status=processing")
			return
		}

		// 3. успех -> обновляем Mongo
		err = s.repo.UpdateStatusAndVariants(
			context.Background(),
			media.ID,
			entity.StatusReady,
			variants,
		)

		if err != nil {
			log.WithError(err).Error("failed to update mongo after processing")
			return
		}

		log.Info("processAsync completed successfully")

	}()
}

func (s *MediaService) GetByIDs(
	ctx context.Context,
	ids []primitive.ObjectID,
) (map[string]entity.MediaDTO, error) {

	if len(ids) == 0 {
		return map[string]entity.MediaDTO{}, nil
	}

	// 1. batch fetch из Mongo
	mediaList, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make(map[string]entity.MediaDTO, len(mediaList))

	for _, m := range mediaList {
		result[m.ID.Hex()] = s.buildDTO(m)
	}

	return result, nil
}

func (s *MediaService) GetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (entity.MediaDTO, error) {

	result, err := s.GetByIDs(ctx, []primitive.ObjectID{id})
	if err != nil {
		return entity.MediaDTO{}, err
	}

	dto, ok := result[id.Hex()]
	if !ok {
		return entity.MediaDTO{}, fmt.Errorf("media not found")
	}

	return dto, nil
}

func (s *MediaService) Delete(ctx context.Context, id primitive.ObjectID) error {

	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	for _, v := range m.Variants {
		err = s.storage.Delete(context.Background(), v.StorageKey)
		if err != nil {
			logrus.WithError(err).WithField("storageKey", v.StorageKey).Warn("failed to delete media variant from storage")
		}
	}

	return nil
}

func (s *MediaService) buildDTO(
	m entity.Media,
) entity.MediaDTO {

	dto := entity.MediaDTO{
		ID:     m.ID.Hex(),
		Status: string(m.Status),
		URLs:   make(map[string]string),
	}

	for _, v := range m.Variants {
		dto.URLs[string(v.Size)] = fmt.Sprintf("/media/%s/%s", m.ID.Hex(), v.Size)
	}

	return dto
}

// HandleStale обрабатывает "зависшие" media, которые долго в статусе processing.
// Если media можно ретраить (retry_count < MaxRetry), то увеличивает retry_count и запускает процессинг заново.
// Если нельзя, то помечает media как failed.
func (s *MediaService) HandleStale(ctx context.Context, limit int) error {

	staleList, err := s.repo.FindStale(ctx, entity.StaleThreshold, limit)
	if err != nil {
		return err
	}

	for _, m := range staleList {

		if !m.CanRetry() {
			_ = s.repo.UpdateStatus(ctx, m.ID, entity.StatusFailed)
			continue
		}

		_ = s.repo.IncrementRetryCount(ctx, m.ID)

		s.processAsyncSizes(m)
	}

	return nil
}
