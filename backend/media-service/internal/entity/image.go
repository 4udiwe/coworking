package entity

import (
	"time"
)

// ──────────────────────────────────────────
// Value Objects
// ──────────────────────────────────────────

type ImageSize string

const (
	SizeThumbnail ImageSize = "thumbnail" // 150px по ширине
	SizeMedium    ImageSize = "medium"    // 800px по ширине
	SizeLarge     ImageSize = "large"     // 1600px по ширине
	SizeOriginal  ImageSize = "original"  // без изменений, конвертируется в WebP
)

// AsyncSizes — размеры которые генерируются асинхронно в фоновой горутине.
// Thumbnail генерируется синхронно и возвращается клиенту сразу в ответе на upload.
var AsyncSizes = []ImageSize{SizeMedium, SizeLarge}

// SizeWidths — ширина в пикселях для каждого размера.
// Height = 0 означает "сохранить пропорции".
var SizeWidths = map[ImageSize]int{
	SizeThumbnail: 150,
	SizeMedium:    800,
	SizeLarge:     1600,
}

type MediaPurpose string

const (
	// PurposeCover — титульная фотография коворкинга. Одна на owner.
	// При загрузке новой — старая автоматически soft-delete'ится.
	PurposeCover MediaPurpose = "cover"
	// PurposeGallery — дополнительные фото для листания, до 20 штук.
	PurposeGallery MediaPurpose = "gallery"
)

type ProcessingStatus string

const (
	// StatusPending — документ создан, файл ещё загружается.
	StatusPending ProcessingStatus = "pending"
	// StatusProcessing — original и thumbnail готовы, medium/large генерируются в фоне.
	StatusProcessing ProcessingStatus = "processing"
	// StatusReady — все варианты готовы, media полностью доступна.
	StatusReady ProcessingStatus = "ready"
	// StatusFailed — обработка не удалась после всех попыток stale checker'а.
	StatusFailed ProcessingStatus = "failed"
)

const (
	OwnerTypeCoworking = "coworking"

	MaxGallerySize = 20               // максимум фото в галерее одного owner
	MaxRetryCount  = 3                // stale checker сдаётся после 3 попыток
	StaleThreshold = 10 * time.Minute // считаем зависшим если processing > 10 мин
)

// ImageVariant — один конкретный размер изображения, хранящийся в MinIO.
type ImageVariant struct {
	Size       ImageSize `bson:"size"`
	Width      int       `bson:"width"`
	Height     int       `bson:"height"`
	StorageKey string    `bson:"storage_key"` // путь внутри bucket: "coworkings/{ownerID}/{mediaID}/thumbnail.webp"
	Bytes      int64     `bson:"bytes"`
}
