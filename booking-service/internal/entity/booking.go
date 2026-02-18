package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Place        Place
	StartTime    time.Time
	EndTime      time.Time
	Status       BookingStatus
	CancelReason *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CancelledAt  *time.Time
}
