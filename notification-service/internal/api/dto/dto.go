package dto

import (
	"time"

	"github.com/google/uuid"
)

// Base DTOs for responses

type NotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Payload   []byte     `json:"payload"`
	IsRead    bool       `json:"isRead"`
	CreatedAt time.Time  `json:"createdAt"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
}

// Request DTOs
type RegisterDeviceRequest struct {
	DeviceToken string `json:"deviceToken" validate:"required"`
	Platform    string `json:"platform" validate:"required"`
}

type MarkNotificationReadRequest struct {
	NotificationID uuid.UUID `json:"notificationId" validate:"required"`
}
