package notification_repository

import (
	"time"

	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
)

type rawNotification struct {
	ID         uuid.UUID  `db:"id"`
	UserID     uuid.UUID  `db:"user_id"`
	TypeID     int16      `db:"notification_type_id"`
	Type       string     `db:"notification_type_name"`
	Title      string     `db:"title"`
	Body       string     `db:"body"`
	Payload    []byte     `db:"payload"`
	StatusID   int16      `db:"status_id"`
	StatusName string     `db:"status_name"`
	CreatedAt  time.Time  `db:"created_at"`
	ReadAt     *time.Time `db:"read_at"`
}

func (r rawNotification) toEntity() entity.Notification {

	return entity.Notification{
		ID:        r.ID,
		UserID:    r.UserID,
		Type:      entity.NotificationType(r.Type),
		Title:     r.Title,
		Body:      r.Body,
		Payload:   r.Payload,
		IsRead:    r.StatusName == string(entity.StatusRead),
		CreatedAt: r.CreatedAt,
		ReadAt:    r.ReadAt,
	}
}
