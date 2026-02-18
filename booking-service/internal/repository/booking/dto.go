package booking_repository

import (
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type rawBookingPlaceStatus struct {
	ID               uuid.UUID  `db:"id"`
	UserID           uuid.UUID  `db:"user_id"`
	PlaceID          uuid.UUID  `db:"place_id"`
	PlaceLabel       string     `db:"place_label"`
	PlaceType        string     `db:"place_type"`
	PlaceCoworkingID uuid.UUID  `db:"place_coworking_id"`
	PlaceIsActive    bool       `db:"place_is_active"`
	StartTime        time.Time  `db:"start_time"`
	EndTime          time.Time  `db:"end_time"`
	StatusID         int        `db:"status_id"`
	StatusName       string     `db:"status_name"`
	CancelReason     *string    `db:"cancel_reason"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	CancelledAt      *time.Time `db:"cancelled_at"`
}

func (r *rawBookingPlaceStatus) toEntity() entity.Booking {
	return entity.Booking{
		ID:     r.ID,
		UserID: r.UserID,
		Place: entity.Place{
			ID:        r.PlaceID,
			Label:     r.PlaceLabel,
			PlaceType: r.PlaceType,
			IsActive:  r.PlaceIsActive,
			Coworking: entity.Coworking{
				ID: r.PlaceCoworkingID,
			},
		},
		StartTime:    r.StartTime,
		EndTime:      r.EndTime,
		Status:       entity.BookingStatus(r.StatusName),
		CancelReason: r.CancelReason,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		CancelledAt:  r.CancelledAt,
	}
}
