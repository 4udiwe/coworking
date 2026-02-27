package post_booking

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

type Request = dto.CreateBookingRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	booking := entity.Booking{
		UserID: in.UserID,
		Place: entity.Place{
			ID: in.PlaceID,
			Coworking: entity.Coworking{
				ID: in.CoworkingID,
			},
		},
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
	}

	err := h.s.CreateBooking(ctx.Request().Context(), booking)

	if err != nil {
		if errors.Is(err, booking_service.ErrBookingStartTimeAfterEndTime) ||
			errors.Is(err, booking_service.ErrBookingStartTimeEqualEndTime) ||
			errors.Is(err, booking_service.ErrBookingStartTimeInPast) ||
			errors.Is(err, booking_service.ErrBookingTimeNotMultipleOfHour) ||
			errors.Is(err, booking_service.ErrBookingDurationLessThanOneHour) ||
			errors.Is(err, booking_service.ErrBookingDurationMoreThanThreeHours) ||
			errors.Is(err, booking_service.ErrPlaceInactive) ||
			errors.Is(err, booking_service.ErrPlaceNotFound) ||
			errors.Is(err, booking_service.ErrCoworkingNotFound) ||
			errors.Is(err, booking_service.ErrCoworkingInactive) ||
			errors.Is(err, booking_service.ErrBookingTimeConflict) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}
