package consumer

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

// Все типы событий, которые могут потребляться сервисом
const (
	BookingCreated   EventType = "booking.created"
	BookingCancelled EventType = "booking.cancelled"
	BookingCompleted EventType = "booking.completed"

	ReminderTriggered EventType = "reminder.triggerred"

	NotificationCreated EventType = "notification.created"
)

// Тип для обработки входящего события
type IncomingEvent struct {
	Type       EventType
	OccurredAt time.Time
	Payload    Payload
}

// Тип данных входящего события. Перечисленны все поля, которые могут быть в событиии.
// (omitempty опускает поле, если его нет)
type Payload struct {
	BookingID        uuid.UUID `json:"bookingId,omitempty"`
	NotificationID   uuid.UUID `json:"notificataionId,omitempty"`
	NotificationType string    `json:"notificataionType,omitempty"`
	UserID           uuid.UUID `json:"userId,omitempty"`
	PlaceID          uuid.UUID `json:"placeId,omitempty"`
	StartTime        time.Time `json:"startTime,omitzero"`
	EndTime          time.Time `json:"endTime,omitzero"`
	Reason           string    `json:"reason,omitempty"`
}
