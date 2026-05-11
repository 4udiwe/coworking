package media_service

import (
	"context"
	"time"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ObjectStorage interface {
	// Upload загружает файл по ключу.
	// data — содержимое файла в памяти (удобно для маленьких файлов).
	// contentType — MIME тип (например "image/webp").
	Upload(ctx context.Context, key string, data []byte, contentType string) error
	// Get скачивает файл по ключу в память.
	// Используется для скачивания original перед async resize.
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete удаляет файл по ключу.
	Delete(ctx context.Context, key string) error
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

	// UpdateStatusAndVariants атомарно обновляет статус и добавляет список вариантов.
	// Используется когда async resize завершён — сразу всё в одном write.
	UpdateStatusAndVariants(ctx context.Context, id primitive.ObjectID, status entity.ProcessingStatus, variants []entity.ImageVariant) error

	Delete(ctx context.Context, id primitive.ObjectID) error

	// Методы для stale checker-а:
	// FindStale возвращает список медиа, которые в статусе Processing уже больше threshold времени.
	FindStale(ctx context.Context, threshold time.Duration, limit int) ([]entity.Media, error)
	// IncrementRetryCount увеличивает на единицу счётчик попыток обработки. Вызывается для stale media.
	IncrementRetryCount(ctx context.Context, id primitive.ObjectID) error
	// UpdateStatus обновляет статус медиа. Вызывается для stale media.
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status entity.ProcessingStatus) error
}
