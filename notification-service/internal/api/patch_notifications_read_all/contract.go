package patch_notifications_read_all

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
}
