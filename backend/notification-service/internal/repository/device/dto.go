package device_repository

import (
	"time"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
)

type rawDevice struct {
	ID uuid.UUID `db:"id"`

	UserID uuid.UUID `db:"user_id"`

	DeviceToken string `db:"device_token"`
	Platform    string `db:"platform"`

	CreatedAt time.Time `db:"created_at"`
}

func (r rawDevice) toEntity() entity.UserDevice {

	return entity.UserDevice{
		ID: r.ID,

		UserID: r.UserID,

		DeviceToken: r.DeviceToken,
		Platform:    r.Platform,

		CreatedAt: r.CreatedAt,
	}
}
