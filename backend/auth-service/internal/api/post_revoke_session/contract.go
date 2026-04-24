package post_revoke_session

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
}
