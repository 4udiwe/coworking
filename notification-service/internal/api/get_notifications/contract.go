package get_notifications

import (
	"context"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
)

type NotificationService interface {
	FetchUnreadNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error)
}
