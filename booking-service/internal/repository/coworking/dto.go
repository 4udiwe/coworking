package coworking_repository

import (
	"time"

	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type rawCoworking struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Address   string    `db:"address"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *rawCoworking) toEntity() *entity.Coworking {
	return &entity.Coworking{
		ID:        r.ID,
		Name:      r.Name,
		Address:   r.Address,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt}
}

type rawCoworkingLayout struct {
	ID          uuid.UUID `db:"id"`
	CoworkingID uuid.UUID `db:"coworking_id"`
	Layout      []byte    `db:"layout"`
	Version     int       `db:"version"`
	CreatedAt   time.Time `db:"created_at"`
}

func (r *rawCoworkingLayout) toEntity() *entity.CoworkingLayout {
	return &entity.CoworkingLayout{
		ID:          r.ID,
		CoworkingID: r.CoworkingID,
		Layout:      r.Layout,
		Version:     r.Version,
		CreatedAt:   r.CreatedAt,
	}
}

type rawLayoutVersionTime struct {
	Version   int       `db:"version"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *rawLayoutVersionTime) toEntity() *entity.CoworkingLayoutVersionTime {
	return &entity.CoworkingLayoutVersionTime{
		Version:   r.Version,
		CreatedAt: r.CreatedAt,
	}
}
