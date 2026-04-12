package get_unread_count

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}
