package get_layout_versions

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	ListLayoutVersions(ctx context.Context, coworkingID uuid.UUID) ([]entity.CoworkingLayoutVersionTime, error)
}
