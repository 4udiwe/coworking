package place_repository

import (
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	"github.com/google/uuid"
)

type rawPlaceCoworking struct {
	ID                uuid.UUID `db:"id"`
	Label             string    `db:"label"`
	PlaceType         string    `db:"place_type"`
	IsActive          bool      `db:"is_active"`
	CoworkingID       uuid.UUID `db:"coworking_id"`
	CoworkingName     string    `db:"coworking_name"`
	CoworkingAddress  string    `db:"coworking_address"`
	CoworkingIsActive bool      `db:"coworking_is_active"`
}

func (r *rawPlaceCoworking) toEntity() entity.Place {
	return entity.Place{
		ID:        r.ID,
		Label:     r.Label,
		PlaceType: r.PlaceType,
		IsActive:  r.IsActive,
		Coworking: entity.Coworking{
			ID:       r.CoworkingID,
			Name:     r.CoworkingName,
			Address:  r.CoworkingAddress,
			IsActive: r.CoworkingIsActive,
		},
	}
}
