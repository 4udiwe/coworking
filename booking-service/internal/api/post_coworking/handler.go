package post_coworking

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
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

type Request struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	coworking := entity.Coworking{
		Name:    in.Name,
		Address: in.Address,
	}
	err := h.s.CreateCoworking(ctx.Request().Context(), coworking)

	if err != nil {
		if errors.Is(err, booking_service.ErrCoworkingAlreadyExists) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}
