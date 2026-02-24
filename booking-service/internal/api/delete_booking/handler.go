package delete_booking

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/4udiwe/coworking/auth-service/pgk/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s BookingService
}

func New(bookingService BookingService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: bookingService})
}

type Request struct {
	BookingID uuid.UUID `json:"bookingId"`
	Reason    string    `json:"reason,omitempty"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	err := h.s.CancelBooking(ctx.Request().Context(), in.BookingID, &in.Reason)

	if err != nil {
		if errors.Is(err, booking_service.ErrBookingNotFound) ||
			errors.Is(err, booking_service.ErrBookingAlreadyCancelled) ||
			errors.Is(err, booking_service.ErrBookingAlreadyCompleted) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusAccepted)
}
