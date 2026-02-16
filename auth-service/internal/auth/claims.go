package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Access token claims
type AccessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

// Refresh token claims
type RefreshClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"` // == refresh_tokens.id
	jwt.RegisteredClaims
}
