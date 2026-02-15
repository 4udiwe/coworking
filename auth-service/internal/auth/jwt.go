package auth

import (
	"errors"
	"time"

	"crypto/sha256"
	"encoding/hex"

	"github.com/samber/lo"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredToken        = errors.New("token has expired")
)

type Auth struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
}

func New(accessTokenSecret, refreshTokenSecret string, accessTokenTTL, refreshTokenTTL time.Duration) *Auth {
	return &Auth{
		accessTokenSecret:  []byte(accessTokenSecret),
		refreshTokenSecret: []byte(refreshTokenSecret),
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
	}
}

func (a *Auth) GenerateTokens(user entity.User) (*Tokens, error) {
	// Access token (15 min TTL)
	accessClaims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  lo.Map(user.Roles, func(r entity.Role, _ int) string { return string(r.Code) }),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(a.accessTokenSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token (7 days)
	refreshClaims := jwt.RegisteredClaims{
		Subject:   user.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(a.refreshTokenSecret)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(a.accessTokenTTL.Seconds()),
	}, nil
}

func (a *Auth) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.accessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidAccessToken
}

func (a *Auth) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.refreshTokenSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", ErrInvalidRefreshToken
}

func (a *Auth) HashToken(tokenString string) string {
	sum := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(sum[:])
}
