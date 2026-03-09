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
	BookingID uuid.UUID `json:"bookingId"`
	UserID    uuid.UUID `json:"userId,omitempty"`
	PlaceID   uuid.UUID `json:"placeId,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Reason    string    `json:"reason,omitempty"`
}
