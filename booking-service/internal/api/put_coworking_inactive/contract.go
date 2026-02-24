package put_coworking_inactive

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	SetCoworkingInActive(ctx context.Context, coworkingID uuid.UUID) error
}
