package dto

import (
	"time"

	"github.com/google/uuid"
)

// Base DTOs for responses

type NotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
}

type UnreadCountResponse struct {
	UnreadCount int `json:"unreadCount"`
}

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Payload   []byte     `json:"payload"`
	ActionURL *string    `json:"actionUrl,omitempty"`
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
	NotificationID uuid.UUID `param:"notificationId" validate:"required"`
}

// Query parameters for fetching notifications
type GetNotificationsQuery struct {
	Limit  int       `query:"limit" validate:"required,min=1,max=100"`
	Offset int       `query:"offset" validate:"min=0"`
	IsRead *bool     `query:"isRead"`
	Since  *time.Time `query:"since"`
}
