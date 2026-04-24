package get_booking_by_id

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	GetBookingByID(ctx context.Context, bookingID uuid.UUID) (entity.Booking, error)
}
