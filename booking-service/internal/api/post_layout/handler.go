package post_layout

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
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

type Request = dto.CreateLayoutRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	layout := entity.CoworkingLayout{
		CoworkingID: in.CoworkingID,
		Version:     in.Version,
		Layout:      in.Layout,
	}
	err := h.s.CreateLayoutVersion(ctx.Request().Context(), layout)

	if err != nil {
		if errors.Is(err, booking_service.ErrCoworkingNotFound) ||
			errors.Is(err, booking_service.ErrInvalidLayoutSchema) ||
			errors.Is(err, booking_service.ErrInvalidLayoutSchemaVersion) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}
