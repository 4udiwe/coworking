package media_service

import (
	"context"
	"time"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectStorage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GeneratePresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}

type ImageProcessor interface {
	ResizeToWidth(ctx context.Context, input []byte, width int) (output []byte, w int, h int, err error)
}

type MediaRepository interface {
	// Create сохраняет новый документ и возвращает его ID.
	Create(ctx context.Context, media entity.Media) (primitive.ObjectID, error)

	// GetByID возвращает медиа по ID. Возвращает ErrMediaNotFound если не найдено.
	GetByID(ctx context.Context, id primitive.ObjectID) (entity.Media, error)

	GetByIDs(ctx context.Context, ids []primitive.ObjectID) ([]entity.Media, error)

	// GetByOwner возвращает все активные (не удалённые) медиа для owner,
	// отсортированные по sort_order.
	GetByOwner(ctx context.Context, ownerType, ownerID string) ([]entity.Media, error)

	// GetCoverByOwner возвращает обложку для owner.
	// Возвращает ErrMediaNotFound если обложки нет.
	GetCoverByOwner(ctx context.Context, ownerType, ownerID string) (entity.Media, error)

	// CountGalleryByOwner возвращает количество фото галереи (не удалённых).
	// Используется для проверки лимита перед загрузкой.
	CountGalleryByOwner(ctx context.Context, ownerType, ownerID string) (int, error)

	// AddVariant атомарно добавляет вариант изображения к документу.
	// Используется после успешного ресайза.
	AddVariant(ctx context.Context, id primitive.ObjectID, variant entity.ImageVariant) error

	// UpdateStatus обновляет статус обработки документа.
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status entity.ProcessingStatus) error

	// UpdateStatusAndVariants атомарно обновляет статус и добавляет список вариантов.
	// Используется когда async resize завершён — сразу всё в одном write.
	UpdateStatusAndVariants(ctx context.Context, id primitive.ObjectID, status entity.ProcessingStatus, variants []entity.ImageVariant) error

	// SoftDelete выставляет deleted_at = now(). Файлы в MinIO не трогает.
	SoftDelete(ctx context.Context, id primitive.ObjectID) error

	// SoftDeleteCoverByOwner мягко удаляет текущую обложку owner'а.
	// Вызывается перед загрузкой новой обложки.
	SoftDeleteCoverByOwner(ctx context.Context, ownerType, ownerID string) error

	// UpdateSortOrder обновляет sort_order для списка медиа.
	// Принимает map[mediaID]newOrder.
	UpdateSortOrder(ctx context.Context, orders map[primitive.ObjectID]int) error

	UpdatePurpose(ctx context.Context, id primitive.ObjectID, purpose entity.MediaPurpose) error
}
