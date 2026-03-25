package put_user_roles

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles []string) error
}
