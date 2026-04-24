package put_coworking

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

type BookingService interface {
	UpdateCoworking(ctx context.Context, coworking entity.Coworking) error
}
