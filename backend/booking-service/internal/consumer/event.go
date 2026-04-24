package consumer

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

// Все типы событий, которые могут потребляться сервисом
const (
	BookingExpire EventType = "booking.expire"
)

// Тип для обработки входящего события
type IncomingEvent struct {
	Type       EventType
	OccurredAt time.Time
	Payload    Payload
}

// Тип данных входящего события. Перечисленны все поля, которые могут быть в событиию.
// (omitempty опускает поле, если его нет)
type Payload struct {
	BookingID uuid.UUID `json:"bookingId"`
}
