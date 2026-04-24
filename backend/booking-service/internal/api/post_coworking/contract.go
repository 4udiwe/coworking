package post_coworking

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

type BookingService interface {
	CreateCoworking(ctx context.Context, coworking entity.Coworking) error
}
