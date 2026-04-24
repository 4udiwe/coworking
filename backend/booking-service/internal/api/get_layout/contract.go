package get_layout

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	GetActiveLayout(ctx context.Context, coworkingID uuid.UUID) (entity.CoworkingLayout, error)
}
