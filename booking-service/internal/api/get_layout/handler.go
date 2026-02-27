package get_layout

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

type Request = dto.GetLayoutRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	layout, err := h.s.GetLatestLayout(ctx.Request().Context(), in.CoworkingID)

	if err != nil {
		if errors.Is(err, booking_service.ErrCoworkingNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, dto.Layout{
		ID:          layout.CoworkingID,
		CoworkingID: layout.CoworkingID,
		Version:     layout.Version,
		Layout:      layout.Layout,
		CreatedAt:   layout.CreatedAt,
	})
}
