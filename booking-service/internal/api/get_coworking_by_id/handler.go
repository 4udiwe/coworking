package get_coworking_by_id

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
	CowokingID uuid.UUID `json:"coworkingId"`
}

type ResponseCoworking struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	c, err := h.s.GetCoworking(ctx.Request().Context(), in.CowokingID)

	if err != nil {
		if errors.Is(err, booking_service.ErrCoworkingNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, ResponseCoworking{
		ID:      c.ID,
		Name:    c.Name,
		Address: c.Address,
	})
}
