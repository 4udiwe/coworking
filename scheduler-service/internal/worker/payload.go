package worker

import "github.com/google/uuid"

// ReminderPayload для события scheduler.reminder.triggered
type ReminderPayload struct {
	BookingID uuid.UUID `json:"bookingId"`
	UserID    uuid.UUID `json:"userId"`
}

// ExpirePayload для события scheduler.booking.expire
type ExpirePayload struct {
	BookingID uuid.UUID `json:"bookingId"`
}