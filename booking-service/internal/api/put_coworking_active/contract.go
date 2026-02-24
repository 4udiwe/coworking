package put_coworking_active

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	SetCoworkingActive(ctx context.Context, coworkingID uuid.UUID) error
}
