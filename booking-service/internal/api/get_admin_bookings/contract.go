package get_admin_bookings

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	ListActiveBookingsForAdmin(ctx context.Context, coworkingID uuid.UUID, page int, pageSize int) ([]entity.Booking, int, error)
}
