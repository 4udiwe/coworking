package post_media

import (
	"context"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
)

type MediaService interface {
	Upload(ctx context.Context, input entity.UploadInput) (entity.UploadResult, error)
}
