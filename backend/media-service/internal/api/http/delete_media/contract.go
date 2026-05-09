package delete_media

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaService interface {
	Delete(ctx context.Context, id primitive.ObjectID) error
}
