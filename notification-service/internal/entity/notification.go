package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationStatus string

const (
	StatusRead   NotificationStatus = "read"
	StatusUnread NotificationStatus = "unread"
)

type NotificationType string

const (
	BookingCreatedNotificationType   NotificationType = "booking_created"
	BookingCancelledNotificationType NotificationType = "booking_cancelled"
	BookingReminderNotificationType  NotificationType = "booking_reminder"
	BookingExpiredNotificationType   NotificationType = "booking_expired"
)

type Notification struct {
	ID uuid.UUID

	UserID uuid.UUID

	Type NotificationType

	Title string
	Body  string

	Payload []byte

	IsRead bool

	CreatedAt time.Time
	ReadAt    *time.Time
}
