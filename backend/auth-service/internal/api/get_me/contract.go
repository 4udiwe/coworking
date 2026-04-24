package get_me

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserInfo(ctx context.Context, userID uuid.UUID) (entity.User, error)
}
