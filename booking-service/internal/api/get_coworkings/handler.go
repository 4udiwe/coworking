package get_coworkings

import (
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/4udiwe/coworking/auth-service/pgk/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s BookingService
}

func New(bookingService BookingService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: bookingService})
}

type Request struct{}

type Response struct {
	Coworkings []ResponseCoworking
}

type ResponseCoworking struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	IsActive bool      `json:"isActive"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	coworkings, err := h.s.ListCoworkings(ctx.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{
		Coworkings: lo.Map(coworkings, func(c entity.Coworking, _ int) ResponseCoworking {
			return ResponseCoworking{
				ID:      c.ID,
				Name:    c.Name,
				Address: c.Address,
				IsActive: c.IsActive,
			}
		}),
	})
}
