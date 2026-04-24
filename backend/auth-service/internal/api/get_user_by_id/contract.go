package get_user_by_id

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserInfo(ctx context.Context, id uuid.UUID) (entity.User, error)
}
