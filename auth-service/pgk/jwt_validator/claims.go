package jwt_validator

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}
