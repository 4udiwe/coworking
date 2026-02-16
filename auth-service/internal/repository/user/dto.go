package user_repository

import (
	"time"

	"github.com/google/uuid"
)

type rawUserRole struct {
	ID            uuid.UUID `db:"id"`
	Email         string    `db:"email"`
	PasswordHash  string    `db:"password_hash"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	RoleID        uuid.UUID `db:"role_id"`
	RoleCode      string    `db:"role_code"`
	RoleName      string    `db:"role_name"`
	RoleCreatedAt time.Time `db:"role_created_at"`
}
