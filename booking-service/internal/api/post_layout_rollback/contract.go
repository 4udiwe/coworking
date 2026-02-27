package post_layout_rollback

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	RollbackLatestLayoutVersion(ctx context.Context, coworkingID uuid.UUID) error
}

