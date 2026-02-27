package put_place_active

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/4udiwe/coworking/auth-service/pgk/decorator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s BookingService
}

func New(bookingService BookingService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: bookingService})
}

type Request = dto.SetPlaceActiveRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	err := h.s.SetPlaceActive(ctx.Request().Context(), in.PlaceID, in.Active)

	if err != nil {
		if errors.Is(err, booking_service.ErrPlaceHasActiveBookings) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if errors.Is(err, booking_service.ErrPlaceNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusAccepted)
}
