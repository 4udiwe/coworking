package get_available_places_by_coworking

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s BookingService
}

func New(bookingService BookingService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: bookingService})
}

type Request = dto.GetAvailablePlacesByCoworkingRequest

type Response struct {
	Places []dto.Place `json:"places"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	places, err := h.s.GetAvailablePlacesByCoworking(ctx.Request().Context(), in.CoworkingID, in.StartTime, in.EndTime)
	if err != nil {
		if errors.Is(err, booking_service.ErrCoworkingNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, Response{
		Places: lo.Map(places, func(p entity.Place, _ int) dto.Place {
			return dto.Place{
				ID:          p.ID,
				CoworkingID: p.Coworking.ID,
				Label:       p.Label,
				PlaceType:   p.PlaceType,
				IsActive:    p.IsActive,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			}
		}),
	})
}
