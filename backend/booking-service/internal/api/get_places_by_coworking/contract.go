package get_places_by_coworking

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	GetPlacesByCoworking(ctx context.Context, coworkingID uuid.UUID) ([]entity.Place, error)
}
