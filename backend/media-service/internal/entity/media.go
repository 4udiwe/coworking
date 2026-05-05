package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Media — основная сущность сервиса.
// Представляет одно загруженное изображение со всеми его вариантами.
type Media struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerType    string             `bson:"owner_type"` // "coworking" — расширяемо на будущие типы
	OwnerID      string             `bson:"owner_id"`   // UUID коворкинга (string для гибкости)
	Purpose      MediaPurpose       `bson:"purpose"`
	OriginalName string             `bson:"original_name"` // оригинальное имя файла от клиента
	MimeType     string             `bson:"mime_type"`     // исходный mime: image/jpeg, image/png, etc.
	Variants     []ImageVariant     `bson:"variants"`
	Status       ProcessingStatus   `bson:"status"`
	SortOrder    int                `bson:"sort_order"`  // порядок в галерее, для cover всегда 0
	UploadedBy   string             `bson:"uploaded_by"` // user_id админа из JWT
	RetryCount   int                `bson:"retry_count"` // сколько раз stale checker перезапускал resize
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"` // nil = активна
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

// IsDeleted проверяет мягкое удаление.
func (m *Media) IsDeleted() bool {
	return m.DeletedAt != nil
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

// StorageKeyFor возвращает путь в MinIO для конкретного размера.
// Формат: "coworkings/{ownerID}/{mediaID}/{size}.webp"
func (m *Media) StorageKeyFor(size ImageSize) string {
	return m.OwnerType + "s/" + m.OwnerID + "/" + m.ID.Hex() + "/" + string(size) + ".webp"
}
