package post_layout

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

type BookingService interface {
	CreateLayoutVersion(ctx context.Context, layout entity.CoworkingLayout) error
}
