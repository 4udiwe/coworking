package get_user_by_id

import (
	"errors"
	"net/http"

	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/dto"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	auth_service "github.com/4udiwe/coworking/auth-service/internal/service/auth"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s UserService
}

func New(s UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: s})
}

type Request = dto.UserByIDRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	user, err := h.s.GetUserInfo(ctx.Request().Context(), in.UserID)

	if err != nil {
		if errors.Is(err, auth_service.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dto.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Roles: lo.Map(user.Roles, func(r entity.Role, _ int) dto.Role {
			return dto.Role{
				ID:       r.ID,
				RoleCode: string(r.Code),
				Name:     r.Name,
			}
		}),
	})
}
