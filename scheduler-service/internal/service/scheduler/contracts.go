package scheduler_service

import (
	"context"

	"github.com/4udiwe/cowoking/scheduler-service/internal/entity"
	"github.com/google/uuid"
)

type TimerRepository interface {
	Create(ctx context.Context, timer entity.Timer) (uuid.UUID, error)
	CancelByBooking(ctx context.Context, bookingID uuid.UUID) error
}
