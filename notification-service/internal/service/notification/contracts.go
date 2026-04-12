package notification_service

import (
	"context"
	"time"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification entity.Notification) (uuid.UUID, error)
	MarkRead(ctx context.Context, id uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Notification, error)
	FetchUnreadByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error)
	FetchByUser(ctx context.Context, userID uuid.UUID, limit, offset int, isRead *bool) ([]entity.Notification, error)
	FetchAfterDate(ctx context.Context, userID uuid.UUID, since time.Time, limit, offset int) ([]entity.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type DeviceRepository interface {
	Create(ctx context.Context, device entity.UserDevice) (uuid.UUID, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserDevice, error)
	DeleteByToken(ctx context.Context, token string) error
}

type OutboxRepository interface {
	Create(ctx context.Context, event entity.OutboxEvent) error
}

type PushService interface {
	SendToUser(ctx context.Context, userID uuid.UUID, notification entity.Notification) error
}
