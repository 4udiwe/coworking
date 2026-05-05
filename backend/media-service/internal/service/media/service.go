package media_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
)

type MediaService struct {
	repo      MediaRepository
	storage   ObjectStorage
	processor ImageProcessor
	baseURL   string

	urlTTL time.Duration
}

func New(repo MediaRepository, storage ObjectStorage, baseURL string) *MediaService {
	return &MediaService{
		repo:    repo,
		storage: storage,
		baseURL: baseURL,
		urlTTL:  15 * time.Minute,
	}
}

func (s *MediaService) Upload(ctx context.Context, input entity.UploadInput) (entity.UploadResult, error) {

	// 1. Проверка лимита галереи
	if input.Purpose == entity.PurposeGallery {
		count, err := s.repo.CountGalleryByOwner(ctx, input.OwnerType, input.OwnerID)
		if err != nil {
			return entity.UploadResult{}, err
		}
		if count >= entity.MaxGallerySize {
			return entity.UploadResult{}, fmt.Errorf("gallery limit exceeded")
		}
	}

	// 2. Cover — удалить старый
	if input.Purpose == entity.PurposeCover {
		if err := s.repo.SoftDeleteCoverByOwner(ctx, input.OwnerType, input.OwnerID); err != nil {
			return entity.UploadResult{}, err
		}
	}

	// 3. Создаём media
	media := entity.Media{
		OwnerType:    input.OwnerType,
		OwnerID:      input.OwnerID,
		Purpose:      input.Purpose,
		OriginalName: input.FileName,
		MimeType:     input.ContentType,
		Status:       entity.StatusPending,
		UploadedBy:   input.UploadedBy,
	}

	id, err := s.repo.Create(ctx, media)
	if err != nil {
		return entity.UploadResult{}, err
	}

	media.ID = id

	// 4. Upload original
	originalKey := media.StorageKeyFor(entity.SizeOriginal)

	if err := s.storage.Upload(ctx, originalKey, input.Data, "image/webp"); err != nil {
		return entity.UploadResult{}, err
	}

	// 4.1 sync resize для thumbnail
	thumbData, _, _, err := s.processor.ResizeToWidth(ctx, input.Data, entity.SizeWidths[entity.SizeThumbnail])
	if err != nil {
		return entity.UploadResult{}, err
	}
	thumbKey := media.StorageKeyFor(entity.SizeThumbnail)

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

	// 8. presigned URL
	url, _ := s.storage.GeneratePresignedURL(ctx, thumbKey, s.urlTTL)

	return entity.UploadResult{
		ID:     id.Hex(),
		Status: string(entity.StatusProcessing),
		URLs: map[string]string{
			"thumbnail": url,
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
		originalKey := media.StorageKeyFor(entity.SizeOriginal)

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

				key := media.StorageKeyFor(size)

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

	// 2. semaphore для ограничения параллелизма
	sem := make(chan struct{}, 10)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, rm := range mediaList {

		m := rm

		wg.Add(1)
		go func() {
			defer wg.Done()

			dto := entity.MediaDTO{
				ID:     m.ID.Hex(),
				Status: string(m.Status),
				URLs:   make(map[string]string),
			}

			for _, v := range m.Variants {

				// ограничиваем concurrency
				sem <- struct{}{}

				url, err := s.storage.GeneratePresignedURL(ctx, v.StorageKey, s.urlTTL)

				<-sem

				if err != nil {
					continue
				}

				dto.URLs[string(v.Size)] = url
			}

			mu.Lock()
			result[m.ID.Hex()] = dto
			mu.Unlock()
		}()
	}

	wg.Wait()

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

func (s *MediaService) GetByOwner(
	ctx context.Context,
	ownerType, ownerID string,
) ([]entity.MediaDTO, error) {

	mediaList, err := s.repo.GetByOwner(ctx, ownerType, ownerID)
	if err != nil {
		return nil, err
	}

	return s.buildDTOList(ctx, mediaList), nil
}

func (s *MediaService) GetCoverByOwner(
	ctx context.Context,
	ownerType, ownerID string,
) (entity.MediaDTO, error) {

	m, err := s.repo.GetCoverByOwner(ctx, ownerType, ownerID)
	if err != nil {
		return entity.MediaDTO{}, err
	}

	dto := s.buildDTO(ctx, m)
	return dto, nil
}

func (s *MediaService) Delete(ctx context.Context, id primitive.ObjectID) error {

	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 1. soft delete
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	// 2. удаляем все варианты
	for _, v := range m.Variants {
		_ = s.storage.Delete(context.Background(), v.StorageKey)
	}

	// 3. удаляем original
	originalKey := m.StorageKeyFor(entity.SizeOriginal)
	_ = s.storage.Delete(context.Background(), originalKey)

	return nil
}

func (s *MediaService) Reorder(
	ctx context.Context,
	orders map[primitive.ObjectID]int,
) error {

	err := s.repo.UpdateSortOrder(ctx, orders)
	if err != nil {
		return err
	}

	return nil
}

func (s *MediaService) SetCover(
	ctx context.Context,
	ownerType, ownerID string,
	mediaID primitive.ObjectID,
) error {

	m, err := s.repo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// 1. валидация
	if m.OwnerType != ownerType || m.OwnerID != ownerID {
		return fmt.Errorf("media does not belong to owner")
	}

	// 2. удалить старую обложку
	if err := s.repo.SoftDeleteCoverByOwner(ctx, ownerType, ownerID); err != nil {
		return err
	}

	// 3. обновить текущую media → сделать cover
	return s.repo.UpdatePurpose(ctx, mediaID, entity.PurposeCover)
}

func (s *MediaService) buildDTO(
	ctx context.Context,
	m entity.Media,
) entity.MediaDTO {

	dto := entity.MediaDTO{
		ID:     m.ID.Hex(),
		Status: string(m.Status),
		URLs:   make(map[string]string),
	}

	for _, v := range m.Variants {

		url, err := s.storage.GeneratePresignedURL(ctx, v.StorageKey, s.urlTTL)
		if err != nil {
			continue
		}

		dto.URLs[string(v.Size)] = url
	}

	return dto
}

func (s *MediaService) buildDTOList(
	ctx context.Context,
	list []entity.Media,
) []entity.MediaDTO {

	result := make([]entity.MediaDTO, len(list))

	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for i, m := range list {

		i, m := i, m

		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			result[i] = s.buildDTO(ctx, m)
		}()
	}

	wg.Wait()
	return result
}
