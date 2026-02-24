package get_booking_by_id

import (
	"errors"
	"net/http"
	"time"

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
	BookingID uuid.UUID `json:"bookingId"`
}

type ResponseBooking struct {
	ID           uuid.UUID     `json:"Id"`
	UserID       uuid.UUID     `json:"userId"`
	Place        ResponsePlace `json:"place"`
	StartTime    time.Time     `json:"startTime"`
	EndTime      time.Time     `json:"endTime"`
	Status       string        `json:"status"`
	CancelReason *string       `json:"cancelReason,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	CancelledAt  *time.Time    `json:"cancelledAt,omitempty"`
}

type ResponsePlace struct {
	ID        uuid.UUID         `json:"id"`
	Coworking ResponseCoworking `json:"coworking"`
	Label     string            `json:"label"`
	PlaceType string            `json:"placeType"`
	IsActive  bool              `json:"isActive"`
}

type ResponseCoworking struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	b, err := h.s.GetBookingByID(ctx.Request().Context(), in.BookingID)

	if err != nil {
		if errors.Is(err, booking_service.ErrBookingNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, ResponseBooking{
		ID:     b.ID,
		UserID: b.UserID,
		Place: ResponsePlace{
			ID: b.Place.ID,
			Coworking: ResponseCoworking{
				ID:      b.Place.Coworking.ID,
				Name:    b.Place.Coworking.Name,
				Address: b.Place.Coworking.Address,
			},
			Label:     b.Place.Label,
			PlaceType: b.Place.PlaceType,
			IsActive:  b.Place.IsActive,
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
