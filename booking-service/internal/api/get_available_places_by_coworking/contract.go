package get_available_places_by_coworking

import (
	"context"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	GetAvailablePlacesByCoworking(ctx context.Context, coworkingID uuid.UUID, start time.Time, end time.Time) ([]entity.Place, error)
}

