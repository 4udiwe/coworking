package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Base DTOs for responses

type Coworking struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	IsActive bool      `json:"isActive"`
}

type Place struct {
	ID          uuid.UUID `json:"id"`
	CoworkingID uuid.UUID `json:"coworkingId"`
	Label       string    `json:"label"`
	PlaceType   string    `json:"placeType"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Booking struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"userId"`
	Place        Place      `json:"place"`
	StartTime    time.Time  `json:"startTime"`
	EndTime      time.Time  `json:"endTime"`
	Status       string     `json:"status"`
	CancelReason *string    `json:"cancelReason,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	CancelledAt  *time.Time `json:"cancelledAt,omitempty"`
}

type Layout struct {
	ID          uuid.UUID       `json:"id"`
	CoworkingID uuid.UUID       `json:"coworkingId"`
	Version     int             `json:"version"`
	Layout      json.RawMessage `json:"layout"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type LayoutVersion struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
}

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Request DTOs

type CreateCoworkingRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UpdateCoworkingRequest struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	IsActive bool      `json:"isActive"`
}

type GetCoworkingByIDRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type ListPlacesByCoworkingRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type GetAvailablePlacesByCoworkingRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId"`
	StartTime   time.Time `query:"startTime"`
	EndTime     time.Time `query:"endTime"`
}

type CreatePlacesRequest struct {
	CoworkingID uuid.UUID        `json:"coworkingId"`
	Places      []CreatePlaceDTO `json:"places"`
}

type CreatePlaceDTO struct {
	Label     string `json:"label"`
	PlaceType string `json:"placeType"`
}

type SetPlaceActiveRequest struct {
	PlaceID uuid.UUID `json:"placeId"`
	Active  bool      `json:"active"`
}

type CreateBookingRequest struct {
	UserID      uuid.UUID `json:"userId"`
	PlaceID     uuid.UUID `json:"placeId"`
	CoworkingID uuid.UUID `json:"coworkingId"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

type GetBookingByIDRequest struct {
	BookingID uuid.UUID `json:"bookingId"`
}

type ListBookingsByUserRequest struct {
	UserID uuid.UUID `json:"userId"`
}

type CancelBookingRequest struct {
	BookingID uuid.UUID `json:"bookingId"`
	Reason    string    `json:"reason,omitempty"`
}

type GetLayoutRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type GetLayoutByVersionRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
	Version     int       `json:"version"`
}

type ListLayoutVersionsRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type CreateLayoutRequest struct {
	CoworkingID uuid.UUID       `json:"coworkingId"`
	Version     int             `json:"version"`
	Layout      json.RawMessage `json:"layout"`
}

type SetCoworkingActiveRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type SetCoworkingInactiveRequest struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}
