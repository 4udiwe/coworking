package sender

import (
	"context"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
)

// Dispatcher interface defines the contract for notification delivery channels
type Dispatcher interface {
	Dispatch(ctx context.Context, notification entity.Notification) error
}
