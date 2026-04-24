package get_layout_by_version

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	 GetLayoutByVersion(ctx context.Context, coworkingID uuid.UUID, version int) (entity.CoworkingLayout, error)
}
