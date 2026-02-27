package auth

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"time"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredToken        = errors.New("token has expired")
)

type Auth struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey

	issuer string

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

/*
New создаёт issuer токенов.

privateKey — приватный RSA ключ
issuer     — идентификатор сервиса (например "auth-service")
*/
func New(
	privateKey *rsa.PrivateKey,
	issuer string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Auth {

	return &Auth{
		privateKey:      privateKey,
		PublicKey:       &privateKey.PublicKey,
		issuer:          issuer,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (a *Auth) GenerateTokens(
	user entity.User,
	sessionID uuid.UUID,
) (*Tokens, error) {

	now := time.Now()

	// ================= ACCESS TOKEN =================

	accessClaims := AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles: lo.Map(user.Roles, func(r entity.Role, _ int) string {
			return string(r.Code)
		}),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(), // jti
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)

	accessTokenString, err := accessToken.SignedString(a.privateKey)
	if err != nil {
		return nil, err
	}

	// ================= REFRESH TOKEN =================

	refreshClaims := RefreshClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.issuer,
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        sessionID.String(), // jti = session id
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(a.privateKey)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(a.accessTokenTTL.Seconds()),
	}, nil
}

func (a *Auth) ParseRefreshToken(
	tokenString string,
) (*RefreshClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {

			// Строгая проверка алгоритма
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, ErrInvalidRefreshToken
			}

			return a.PublicKey, nil
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

	// Проверка issuer
	if claims.Issuer != a.issuer {
		return nil, ErrInvalidRefreshToken
	}

	return claims, nil
}

func LoadPrivateKeyFromPEM(pemData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func (a *Auth) HashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}
