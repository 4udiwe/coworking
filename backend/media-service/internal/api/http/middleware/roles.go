package middleware

import (
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type RoleCode string

const (
	RoleStudent RoleCode = "student"
	RoleTeacher RoleCode = "teacher"
	RoleAdmin   RoleCode = "admin"
)

func RoleMiddleware(allowedRoles ...RoleCode) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetUserFromContext(c)
			if err != nil {
				logrus.Errorf("Role middleware: get user from context error:%v", err)
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
	return RoleMiddleware(RoleAdmin)(next)
}

func TeacherOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(RoleTeacher)(next)
}

func AdminAndTeacher(next echo.HandlerFunc) echo.HandlerFunc {
	return RoleMiddleware(RoleAdmin, RoleTeacher)(next)
}
