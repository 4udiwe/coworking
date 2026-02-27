package post_places

import (
	"errors"
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/4udiwe/coworking/auth-service/pgk/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s BookingService
}

func New(bookingService BookingService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: bookingService})
}

type Request = dto.CreatePlacesRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	places := lo.Map(in.Places, func(p dto.CreatePlaceDTO, _ int) entity.Place {
		return entity.Place{
			Coworking: entity.Coworking{ID: in.CoworkingID},
			Label:     p.Label,
			PlaceType: p.PlaceType,
		}
	})
	err := h.s.CreatePlacesBatch(ctx.Request().Context(), places)

	if err != nil {
		if errors.Is(err, booking_service.ErrPlaceAlreadyExists) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}
