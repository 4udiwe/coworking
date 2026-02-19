package booking_service

import (
	"context"
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type BookingRepository interface {
	Create(ctx context.Context, booking entity.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Booking, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error)
	Cancel(ctx context.Context, id uuid.UUID, reason *string) error
	MarkCompleted(ctx context.Context, id uuid.UUID) error
}

type PlaceRepository interface {
	CreateBatch(ctx context.Context, places []entity.Place) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Place, error)
	GetByCoworking(ctx context.Context, coworkingID uuid.UUID) ([]entity.Place, error)
	GetAvailableByCoworking(ctx context.Context, coworkingID uuid.UUID, start time.Time, end time.Time) ([]entity.Place, error)
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
	CheckHasActiveBookings(ctx context.Context, placeID uuid.UUID) (bool, error)
}

type CoworkingRepository interface {
	Create(ctx context.Context, coworking entity.Coworking) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Coworking, error)
	Update(ctx context.Context, coworking entity.Coworking) error
	List(ctx context.Context) ([]entity.Coworking, error)
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
	CheckHasActiveBookings(ctx context.Context, coworkingID uuid.UUID) (bool, error)

	CreateLayoutVersion(ctx context.Context, layout entity.CoworkingLayout) error
	RollbackLatestLayoutVersion(ctx context.Context, coworkingID uuid.UUID) error
	GetLatestLayout(ctx context.Context, coworkingID uuid.UUID) (entity.CoworkingLayout, error)
	GetLayoutByVersion(ctx context.Context, coworkingID uuid.UUID, version int) (entity.CoworkingLayout, error)
	ListLayoutVersions(ctx context.Context, coworkingID uuid.UUID) ([]entity.CoworkingLayoutVersionTime, error)
}
