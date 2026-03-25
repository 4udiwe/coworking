package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Base DTOs for responses

type Coworking struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	UserName     string     `json:"userName"`
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

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// Request DTOs

type CreateCoworkingRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=200"`
	Address string `json:"address" validate:"required,min=5,max=500"`
}

type UpdateCoworkingRequest struct {
	ID       uuid.UUID `param:"coworkingId" validate:"required"`
	Name     string    `json:"name" validate:"required,min=2,max=200"`
	Address  string    `json:"address" validate:"required,min=5,max=500"`
	IsActive bool      `json:"isActive"`
}

type GetCoworkingByIDRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
}

type ListPlacesByCoworkingRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
}

type GetAvailablePlacesByCoworkingRequest struct {
	CoworkingID uuid.UUID  `param:"coworkingId" validate:"required"`
	StartTime   *time.Time `query:"startTime" validate:"required"`
	EndTime     *time.Time `query:"endTime" validate:"required"`
}

type CreatePlacesRequest struct {
	CoworkingID uuid.UUID        `json:"coworkingId" validate:"required"`
	Places      []CreatePlaceDTO `json:"places" validate:"required,min=1,dive"`
}

type CreatePlaceDTO struct {
	Label     string `json:"label" validate:"required,min=1,max=50"`
	PlaceType string `json:"placeType" validate:"required,oneof=open_desk meeting_room private_office"`
}

type SetPlaceActiveRequest struct {
	PlaceID uuid.UUID `param:"placeId" validate:"required"`
	Active  *bool     `json:"active"`
}

type CreateBookingRequest struct {
	PlaceID   uuid.UUID `json:"placeId" validate:"required"`
	StartTime time.Time `json:"startTime" validate:"required"`
	EndTime   time.Time `json:"endTime" validate:"required,gtfield=StartTime"`
}

type GetBookingByIDRequest struct {
	BookingID uuid.UUID `param:"bookingId" validate:"required"`
}

type GetAdminActiveBookingsRequest struct {
	CoworkingID uuid.UUID `query:"coworkingId" validate:"required"`
	Page        int       `query:"page" validate:"required,min=1"`
	PageSize    int       `query:"pageSize" validate:"required,min=1,max=100"`
}

type GetBookingsByUserRequest struct {
	Page       int    `query:"page" validate:"required,min=1"`
	PageSize   int    `query:"pageSize" validate:"required,min=1,max=100"`
	Status     *string `query:"status" validate:"omitempty,oneof=active completed cancelled"`
}

type CancelBookingRequest struct {
	BookingID uuid.UUID `param:"bookingId" validate:"required"`
	Reason    string    `json:"reason,omitempty" validate:"max=500"`
}

type GetLayoutRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
}

type GetLayoutByVersionRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
	Version     int       `param:"version" validate:"required,gt=0"`
}

type ListLayoutVersionsRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
}

type CreateLayoutRequest struct {
	CoworkingID uuid.UUID       `json:"coworkingId" validate:"required"`
	Version     int             `json:"version" validate:"required,gt=0"`
	Layout      json.RawMessage `json:"layout" validate:"required"`
}

type DeleteLayoutRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
	Version     int       `param:"version" validate:"required,gt=0"`
}

type SetActiveLayoutRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
	Version     int       `param:"version" validate:"required,gt=0"`
}

type SetCoworkingActiveRequest struct {
	CoworkingID uuid.UUID `param:"coworkingId" validate:"required"`
	Active      *bool     `json:"active" validate:"required"`
}
