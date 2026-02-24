package get_layout_versions

import (
	"errors"
	"net/http"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
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

type Request struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type Response struct {
	Versions []ResponseVersion `json:"layoutVersions"`
}

type ResponseVersion struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
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
		Versions: lo.Map(versions, func(v entity.CoworkingLayoutVersionTime, _ int) ResponseVersion {
			return ResponseVersion{
				Version:   v.Version,
				CreatedAt: v.CreatedAt,
			}
		}),
	})
}
