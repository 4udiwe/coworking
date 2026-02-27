package put_place_active

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	SetPlaceActive(ctx context.Context, placeID uuid.UUID, active bool) error
}
