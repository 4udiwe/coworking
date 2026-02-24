package get_coworking_by_id

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	GetCoworking(ctx context.Context, coworkingID uuid.UUID) (entity.Coworking, error)
}
