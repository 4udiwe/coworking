package post_device

import (
	"context"

	"github.com/google/uuid"
)

type NotificationService interface {
	RegisterDevice(ctx context.Context, userID uuid.UUID, deviceToken string, platform string) error
}
