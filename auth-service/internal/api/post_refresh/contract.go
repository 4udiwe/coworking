package post_refresh

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
)

type UserService interface {
	Refresh(ctx context.Context, refreshToken string) (*auth.Tokens, error)
}
