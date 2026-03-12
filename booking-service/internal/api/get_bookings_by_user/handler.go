package get_bookings_by_user

import (
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
	"github.com/4udiwe/cowoking/booking-service/internal/api/middleware"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
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

type Request = struct{}

type Response struct {
	Bookings []dto.Booking `json:"bookings"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	bookings, err := h.s.ListBookingsByUser(ctx.Request().Context(), claims.UserID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{
		Bookings: lo.Map(bookings, func(b entity.Booking, _ int) dto.Booking {
			return dto.Booking{
				ID:     b.ID,
				UserID: b.UserID,
				Place: dto.Place{
					ID:          b.Place.ID,
					CoworkingID: b.Place.Coworking.ID,
					Label:       b.Place.Label,
					PlaceType:   b.Place.PlaceType,
					IsActive:    b.Place.IsActive,
					CreatedAt:   b.Place.CreatedAt,
					UpdatedAt:   b.Place.UpdatedAt,
				},
				StartTime:    b.StartTime,
				EndTime:      b.EndTime,
				Status:       string(b.Status),
				CancelReason: b.CancelReason,
				CreatedAt:    b.CreatedAt,
				UpdatedAt:    b.UpdatedAt,
				CancelledAt:  b.CancelledAt,
			}
		}),
	})
}
