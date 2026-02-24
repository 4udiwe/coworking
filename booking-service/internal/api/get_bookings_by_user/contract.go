package get_bookings_by_user

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	ListBookingsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error)
}
