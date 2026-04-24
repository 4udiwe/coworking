package delete_booking

import (
	"context"

	"github.com/google/uuid"
)

type BookingService interface {
	CancelBooking(ctx context.Context, bookingID uuid.UUID, reason *string) error
}
