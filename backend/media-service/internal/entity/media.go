package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Media — основная сущность сервиса.
// Представляет одно загруженное изображение со всеми его вариантами.
type Media struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	OriginalName string             `bson:"original_name"` // оригинальное имя файла от клиента
	MimeType     string             `bson:"mime_type"`     // исходный mime: image/jpeg, image/png, etc.
	Variants     []ImageVariant     `bson:"variants"`
	Status       ProcessingStatus   `bson:"status"`
	UploadedBy   string             `bson:"uploaded_by"` // user_id админа из JWT
	RetryCount   int                `bson:"retry_count"` // сколько раз stale checker перезапускал resize
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	ExpiresAt    *time.Time         `bson:"expires_at,omitempty"` // для soft delete
}

type MediaDTO struct {
	ID     string
	Status string
	URLs   map[string]string
}

// GetVariant возвращает вариант нужного размера если он уже сгенерирован.
func (m *Media) GetVariant(size ImageSize) (ImageVariant, bool) {
	for _, v := range m.Variants {
		if v.Size == size {
			return v, true
		}
	}
	return ImageVariant{}, false
}

// IsStale возвращает true если документ завис в processing дольше StaleThreshold.
// Используется stale checker'ом для поиска потерянных задач.
func (m *Media) IsStale() bool {
	return m.Status == StatusProcessing &&
		time.Since(m.UpdatedAt) > StaleThreshold
}

// CanRetry возвращает true если stale checker ещё может попробовать перезапустить.
func (m *Media) CanRetry() bool {
	return m.RetryCount < MaxRetryCount
}

// NeedsAsyncResize возвращает true если ещё не все async-варианты сгенерированы.
func (m *Media) NeedsAsyncResize() bool {
	for _, size := range AsyncSizes {
		if _, ok := m.GetVariant(size); !ok {
			return true
		}
	}
	return false
}
