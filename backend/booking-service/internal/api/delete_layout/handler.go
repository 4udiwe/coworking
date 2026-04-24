package delete_layout

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

type Request = dto.DeleteLayoutRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	err := h.s.DeleteLayoutVersion(ctx.Request().Context(), in.CoworkingID, in.Version)

	if err != nil {
		if errors.Is(err, booking_service.ErrLayoutNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusAccepted)
}
