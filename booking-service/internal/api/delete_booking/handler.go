package delete_booking

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s BookingService
}

func New(bookingService BookingService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: bookingService})
}

type Request = dto.CancelBookingRequest

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
