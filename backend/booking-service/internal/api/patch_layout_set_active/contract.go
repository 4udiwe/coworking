package patch_layout_set_active

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	SetLayoutVersionToActive(ctx context.Context, coworkingID uuid.UUID, layoutVersion int) error
}
