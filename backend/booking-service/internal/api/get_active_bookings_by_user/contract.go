package get_active_bookings_by_user

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	ListActiveBookingsByUser(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]entity.Booking, int, error)
}
