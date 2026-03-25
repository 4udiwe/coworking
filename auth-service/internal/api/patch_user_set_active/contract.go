package patch_user_set_active

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	SetUserActive(ctx context.Context, userID uuid.UUID, active bool) error
}
