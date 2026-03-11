package notification_service

import (
	"context"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/4udiwe/coworking/notification-service/internal/sender"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification entity.Notification) (uuid.UUID, error)
	MarkRead(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Notification, error)
	FetchUnreadByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error)
}

type DeviceRepository interface {
	Create(ctx context.Context, device entity.UserDevice) (uuid.UUID, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserDevice, error)
	DeleteByToken(ctx context.Context, token string) error
}

type OutboxRepository interface {
	Create(ctx context.Context, event entity.OutboxEvent) error
}

type PushSender interface {
	Send(ctx context.Context, msg sender.PushMessage) error
}
