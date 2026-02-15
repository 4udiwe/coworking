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

type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}
