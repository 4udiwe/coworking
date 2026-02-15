package user_service

import (
	"context"
	"time"

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
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetUserByRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type Auth interface {
	GenerateTokens(user entity.User) (*auth.Tokens, error)
	ValidateAccessToken(tokenString string) (*auth.TokenClaims, error)
	ValidateRefreshToken(tokenString string) (string, error)
	HashToken(tokenString string) string
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}
