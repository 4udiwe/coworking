package user_service

import (
	"context"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mocks.go -package=mocks

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	AttachRole(ctx context.Context, userID uuid.UUID, roleCode string) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type AuthRepository interface {
	CreateSession(ctx context.Context, session entity.Session, tokenHash string) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (entity.Session, error)
	UpdateLastUsedAt(ctx context.Context, id uuid.UUID) error
	RevokeSession(ctx context.Context, id uuid.UUID) error
	GetUserSessions(ctx context.Context, userID uuid.UUID, onlyActive bool) ([]entity.Session, error)
}

type Auth interface {
	GenerateTokens(user entity.User, sessionID uuid.UUID) (*auth.Tokens, error)
	ParseRefreshToken(tokenString string) (*auth.RefreshClaims, error)
	HashToken(tokenString string) string
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}
