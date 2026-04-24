package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	UserAgent         string
	IPAddress         string
	DeviceName        *string
	DeviceFingerprint *string
	ExpiresAt         time.Time
	LastUsedAt        time.Time
	Revoked           bool
	CreatedAt         time.Time
}
