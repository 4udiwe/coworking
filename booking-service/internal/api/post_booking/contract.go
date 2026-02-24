package post_booking

import (
	"context"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
)

type BookingService interface {
	CreateBooking(ctx context.Context, booking entity.Booking) error
}
