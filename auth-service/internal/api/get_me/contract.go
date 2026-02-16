package get_me

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
)

type UserService interface {
	GetUserInfo(ctx context.Context, refreshToken string) (entity.User, error)
}
