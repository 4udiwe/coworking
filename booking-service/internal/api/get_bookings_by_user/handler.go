package get_bookings_by_user

import (
	"net/http"
	"time"

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

type Request struct {
	UserID uuid.UUID `json:"userId"`
}

type Response struct {
	Bookings []ResponseBooking `json:"bookings"`
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
	bookings, err := h.s.ListBookingsByUser(ctx.Request().Context(), in.UserID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{
		Bookings: lo.Map(bookings, func(b entity.Booking, _ int) ResponseBooking {
			return ResponseBooking{
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
			}
		}),
	})
}
