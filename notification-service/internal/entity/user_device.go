package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserDevice struct {
	ID uuid.UUID

	UserID uuid.UUID

	DeviceToken string
	Platform    string

	CreatedAt time.Time
}
