package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookingEvent struct {
	EventID   uuid.UUID
	EventType string

	BookingID   uuid.UUID
	CoworkingID uuid.UUID
	UserID      uuid.UUID
	PlaceID     uuid.UUID

	StartTime time.Time
	EndTime   time.Time

	BookingStatus string

	Occurred time.Time
}
