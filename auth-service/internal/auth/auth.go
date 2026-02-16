package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/samber/lo"
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

func New(
	accessTokenSecret string,
	refreshTokenSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Auth {
	return &Auth{
		accessTokenSecret:  []byte(accessTokenSecret),
		refreshTokenSecret: []byte(refreshTokenSecret),
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
	}
}

func (a *Auth) GenerateTokens(
	user entity.User,
	sessionID uuid.UUID,
) (*Tokens, error) {

	now := time.Now()

	// ===== ACCESS TOKEN =====
	accessClaims := AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles: lo.Map(user.Roles, func(r entity.Role, _ int) string {
			return string(r.Code)
		}),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(a.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessTokenString, err := accessToken.SignedString(a.accessTokenSecret)
	if err != nil {
		return nil, err
	}

	// ===== REFRESH TOKEN =====
	refreshClaims := RefreshClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sessionID.String(), // jti
			ExpiresAt: jwt.NewNumericDate(now.Add(a.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
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

func (a *Auth) ValidateAccessToken(tokenString string) (*AccessClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return a.accessTokenSecret, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidAccessToken
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidAccessToken
	}

	return claims, nil
}

func (a *Auth) ParseRefreshToken(
	tokenString string,
) (*RefreshClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return a.refreshTokenSecret, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	return claims, nil
}

func (a *Auth) HashToken(tokenString string) string {
	sum := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(sum[:])
}
