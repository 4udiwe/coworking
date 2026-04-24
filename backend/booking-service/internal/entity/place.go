package entity

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID        uuid.UUID
	Coworking Coworking
	Label     string
	PlaceType string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
