package get_active_sessions

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
)

type UserService interface {
	GetUserSessions(
		ctx context.Context,
		refreshToken string,
		onlyActive bool,
	) ([]entity.Session, error)
}
