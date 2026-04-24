package worker

import (
	"time"

	"github.com/google/uuid"
)

// ReminderPayload для события scheduler.reminder.triggered
type ReminderPayload struct {
	BookingID  uuid.UUID  `json:"bookingId"`
	UserID     uuid.UUID  `json:"userId"`
	PlaceID    *uuid.UUID `json:"placeId,omitempty"`
	PlaceLabel *string    `json:"placeLabel,omitempty"`
	StartTime  *time.Time `json:"startTime,omitempty"`
	EndTime    *time.Time `json:"endTime,omitempty"`
}

type ExpirePayload struct {
	BookingID uuid.UUID `json:"bookingId"`
}
