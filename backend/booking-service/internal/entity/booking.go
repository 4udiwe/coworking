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
	UserName     string
	Place        Place
	StartTime    time.Time
	EndTime      time.Time
	Status       BookingStatus
	CancelReason *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CancelledAt  *time.Time
}

// Возвращает все статусы с которыми бронирования отображаются в приложении
// на вкладке "Active" в разделе бронирований пользователя
func GetActiveStatuses() []string {
	return []string{string(BookingStatusActive)}
}

// Возвращает все статусы с которыми бронирования отображаются в приложении
// на вкладке "History" в разделе бронирований пользователя
func GetHistoryStatuses() []string {
	return []string{string(BookingStatusCancelled), string(BookingStatusCompleted)}
}
