package post_register

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
)

type UserService interface {
	Register(
		ctx context.Context,
		email string,
		password string,
		firstName string,
		lastName string,
		roleCode string,
		userAgent string,
		deviceInfo string,
		ip string,
	) (*auth.Tokens, error)
}
