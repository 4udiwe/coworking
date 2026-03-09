package entity

import (
	"time"

	"github.com/google/uuid"
)

type TimerID int16

const (
	TimerTypeBookingReminderID TimerID = 1
	TimerTypeBookingExpireID   TimerID = 2
)

type TimerName string

const (
	TimerTypeBookingReminderName TimerName = "booking_reminder"
	TimerTypeBoookingExpireName  TimerName = "booking_expire"
)

type TimerStatus string

const (
	TimerStatusPending   TimerStatus = "pending"
	TimerStatusTriggered TimerStatus = "triggered"
	TimerStatusCancelled TimerStatus = "cancelled"
)

type TimerType struct {
	ID   TimerID
	Name TimerName
}

type Timer struct {
	ID uuid.UUID

	Type TimerType

	BookingID uuid.UUID
	UserID    *uuid.UUID

	TriggerAt time.Time

	Payload []byte

	Status TimerStatus

	CreatedAt   time.Time
	TriggeredAt *time.Time
	CancelledAt *time.Time
}
