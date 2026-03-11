package notification_builder

import (
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
)

type Event struct {
	Type entity.NotificationType

	UserID uuid.UUID

	Payload map[string]any
}
