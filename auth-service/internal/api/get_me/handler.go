package get_me

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/middleware"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request struct{}

type ResponseUser struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	IsActive  bool           `json:"isActive"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Roles     []ResponseRole `json:"roles"`
}

type ResponseRole struct {
	ID       string `json:"id"`
	RoleCode string `json:"roleCode"`
	Name     string `json:"name"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	user, err := h.s.GetUserInfo(ctx.Request().Context(), claims.UserID)

	if err != nil {
		// Validation errors
		if errors.Is(err, user_service.ErrEmptyUserID) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		// Authentication errors
		if errors.Is(err, user_service.ErrUserInactive) ||
			errors.Is(err, user_service.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		// Any other error is internal server error
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, ResponseUser{
		ID:        user.ID.String(),
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Roles: lo.Map(user.Roles, func(role entity.Role, _ int) ResponseRole {
			return ResponseRole{
				ID:       role.ID.String(),
				RoleCode: string(role.Code),
				Name:     role.Name,
			}
		}),
	})
}
