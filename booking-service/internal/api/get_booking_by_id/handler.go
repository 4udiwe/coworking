package get_booking_by_id

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

type Request = dto.GetBookingByIDRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	b, err := h.s.GetBookingByID(ctx.Request().Context(), in.BookingID)

	if err != nil {
		if errors.Is(err, booking_service.ErrBookingNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, dto.Booking{
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
	})
}
