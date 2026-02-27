package middleware

import (
	"net/http"
	"strings"

	"github.com/4udiwe/coworking/auth-service/pgk/jwt_validator"
	"github.com/labstack/echo/v4"
)

const USER_CLAIMS_KEY = "userClaims"

type AuthMiddleware struct {
	jwtValidator *jwt_validator.Validator
}

func New(jwtValidator *jwt_validator.Validator) *AuthMiddleware {
	return &AuthMiddleware{
		jwtValidator: jwtValidator,
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
		claims, err := m.jwtValidator.Validate(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set(USER_CLAIMS_KEY, claims)

		return next(c)
	}
}

func GetUserFromContext(c echo.Context) (*jwt_validator.AccessClaims, error) {
	claims, ok := c.Get(USER_CLAIMS_KEY).(*jwt_validator.AccessClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}
	return claims, nil
}
