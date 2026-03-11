package patch_notification

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	MarkRead(ctx context.Context, notificationID uuid.UUID) error
}
