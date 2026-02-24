package get_coworkings

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

type BookingService interface {
	ListCoworkings(ctx context.Context) ([]entity.Coworking, error)
}
