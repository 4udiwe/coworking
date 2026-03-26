package get_admin_bookings

import (
	"context"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingService interface {
	ListActiveBookingsForAdmin(ctx context.Context, coworkingID uuid.UUID, page int, pageSize int, dateFrom *time.Time, dateTo *time.Time, placeType *string, sortBy *string) ([]entity.Booking, int, error)
}
