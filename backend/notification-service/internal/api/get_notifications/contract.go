package get_notifications

import (
	"context"
	"time"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
)

type NotificationService interface {
	FetchUnreadNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error)
	FetchNotifications(ctx context.Context, userID uuid.UUID, limit, offset int, isRead *bool) ([]entity.Notification, error)
	FetchNotificationsAfterDate(ctx context.Context, userID uuid.UUID, since time.Time, limit, offset int) ([]entity.Notification, error)
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
}
