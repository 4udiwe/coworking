package timer_repository

import (
	"time"

	"github.com/4udiwe/cowoking/scheduler-service/internal/entity"
	"github.com/google/uuid"
)

type rawTimer struct {
	ID uuid.UUID `db:"id"`

	TimerTypeID int16 `db:"timer_type_id"`

	BookingID uuid.UUID  `db:"booking_id"`
	UserID    *uuid.UUID `db:"user_id"`

	TriggerAt time.Time `db:"trigger_at"`

	Payload []byte `db:"payload"`

	StatusID int16 `db:"status_id"`

	CreatedAt   time.Time  `db:"created_at"`
	TriggeredAt *time.Time `db:"triggered_at"`
	CancelledAt *time.Time `db:"cancelled_at"`
}

func (r rawTimer) toEntity() entity.Timer {

	return entity.Timer{
		ID: r.ID,

		BookingID: r.BookingID,
		UserID:    r.UserID,

		Type: entity.TimerType{
			ID: entity.TimerID(r.TimerTypeID),
		},

		TriggerAt: r.TriggerAt,
		Payload:   r.Payload,

		CreatedAt:   r.CreatedAt,
		TriggeredAt: r.TriggeredAt,
		CancelledAt: r.CancelledAt,
	}
}
