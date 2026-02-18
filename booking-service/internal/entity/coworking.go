package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Coworking struct {
	ID        uuid.UUID
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CoworkingLayout struct {
	ID          uuid.UUID
	CoworkingID uuid.UUID
	Layout      json.RawMessage
	Version     int
	CreatedAt   time.Time
}

type CoworkingLayoutVersionTime struct {
	Version   int
	CreatedAt time.Time
}
