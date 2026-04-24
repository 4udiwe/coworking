package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetUsersRequest struct {
	Search   string  `query:"search"`
	Page     *int    `query:"page"`
	Size     *int    `query:"size"`
	Role     *string `query:"role"`
	IsActive *bool   `query:"isActive"`
	Sort     string  `query:"sort"`
}

type UserByIDRequest struct {
	UserID uuid.UUID `param:"userId" validate:"required"`
}

type PaginatedUsers struct {
	Users []User `json:"users"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int64  `json:"total"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Roles     []Role    `json:"roles"`
}

type Role struct {
	ID       uuid.UUID `json:"id"`
	RoleCode string    `json:"roleCode"`
	Name     string    `json:"name"`
}
