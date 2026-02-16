package middleware

import (
	"net/http"
	"strings"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/labstack/echo/v4"
)

const USER_CLAIMS_KEY = "userClaims"

type AuthRepo interface {
	ValidateAccessToken(tokenString string) (*auth.AccessClaims, error)
}

type AuthMiddleware struct {
	auth AuthRepo
}

func New(auth AuthRepo) *AuthMiddleware {
	return &AuthMiddleware{
		auth: auth,
	}
}

func (m *AuthMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
		}

		token := parts[1]
		claims, err := m.auth.ValidateAccessToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set(USER_CLAIMS_KEY, claims)

		return next(c)
	}
}

func GetUserFromContext(c echo.Context) (*auth.AccessClaims, error) {
	claims, ok := c.Get(USER_CLAIMS_KEY).(*auth.AccessClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}
	return claims, nil
}
