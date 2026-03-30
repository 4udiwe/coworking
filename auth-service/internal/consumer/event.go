package consumer

import (
	"time"
)

type EventType string

// Все типы событий, которые потребляет auth-service
const (
	SessionsCleanup EventType = "sessions.cleanup"
)

// Тип для обработки входящего события
type IncomingEvent struct {
	Type       EventType
	OccurredAt time.Time
	Payload    Payload
}

// Payload события очистки сессий
type Payload struct {
	RetentionDays int `json:"retentionDays"`
}
