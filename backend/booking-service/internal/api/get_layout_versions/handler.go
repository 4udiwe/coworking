package get_layout_versions

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

type Request = dto.ListLayoutVersionsRequest

type Response struct {
	Versions []dto.LayoutVersion `json:"layoutVersions"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	versions, err := h.s.ListLayoutVersions(ctx.Request().Context(), in.CoworkingID)

	if err != nil {
		if errors.Is(err, booking_service.ErrCoworkingNotFound) {
			echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{
		Versions: lo.Map(versions, func(v entity.CoworkingLayoutVersionTime, _ int) dto.LayoutVersion {
			return dto.LayoutVersion{
				Version:   v.Version,
				CreatedAt: v.CreatedAt,
			}
		}),
	})
}
