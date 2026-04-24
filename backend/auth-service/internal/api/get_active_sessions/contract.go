package get_active_sessions

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserSessions(
		ctx context.Context,
		userID uuid.UUID,
		onlyActive bool,
	) ([]entity.Session, error)
}
