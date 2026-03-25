package get_admin_bookings

import (
	"net/http"

	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/dto"
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

type Request = dto.GetAdminActiveBookingsRequest

type Response struct {
	Bookings   []dto.Booking       `json:"bookings"`
	Pagination dto.PaginationMeta  `json:"pagination"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	bookings, totalCount, err := h.s.ListActiveBookingsForAdmin(ctx.Request().Context(), in.CoworkingID, in.Page, in.PageSize)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalPages := (totalCount + in.PageSize - 1) / in.PageSize

	return ctx.JSON(http.StatusOK, Response{
		Bookings: lo.Map(bookings, func(b entity.Booking, _ int) dto.Booking {
			return dto.Booking{
				ID:       b.ID,
				UserID:   b.UserID,
				UserName: b.UserName,
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
		Pagination: dto.PaginationMeta{
			Page:       in.Page,
			PageSize:   in.PageSize,
			TotalItems: totalCount,
			TotalPages: totalPages,
		},
	})
}
