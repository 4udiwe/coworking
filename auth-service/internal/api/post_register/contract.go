package post_register

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
)

type UserService interface {
	Register(ctx context.Context, email string, password string, roleCode string) (*auth.Tokens, error)
}
