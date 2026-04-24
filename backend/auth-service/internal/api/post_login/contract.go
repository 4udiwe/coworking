package post_login

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
)

type UserService interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		userAgent string,
		deviceInfo string,
		ip string,
	) (*auth.Tokens, error)
}
