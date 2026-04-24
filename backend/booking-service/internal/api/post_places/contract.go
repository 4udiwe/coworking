package post_places

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

type BookingService interface {
	CreatePlacesBatch(ctx context.Context, places []entity.Place) error
}
