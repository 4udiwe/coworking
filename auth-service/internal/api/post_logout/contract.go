package post_logout

import (
	"context"
)

type UserService interface {
	Logout(
		ctx context.Context,
		refreshToken string,
	) error
}
