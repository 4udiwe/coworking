package middleware

import (
	"net/http"
	"slices"

	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles ...entity.RoleCode) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetUserFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			hasAccess := false
			for _, allowedRole := range allowedRoles {
				if slices.Contains(claims.Roles, string(allowedRole)) {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied")
			}

			return next(c)
		}
	}
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(entity.RoleAdmin)(next)
}

func TeacherOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(entity.RoleTeacher)(next)
}

func AdminAndTeacher(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(entity.RoleAdmin, entity.RoleTeacher)(next)
}
