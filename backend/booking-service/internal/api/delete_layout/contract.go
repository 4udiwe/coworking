package delete_layout

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	DeleteLayoutVersion(ctx context.Context, coworkingID uuid.UUID, layoutVersion int) error
}
