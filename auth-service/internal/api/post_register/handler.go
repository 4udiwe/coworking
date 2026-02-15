package post_register

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/decorator"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	RoleCode string `json:"roleCode" validate:"required"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	tokens, err := h.s.Register(ctx.Request().Context(), in.Email, in.Password, in.RoleCode)

	if err != nil {
		// Validation errors
		if errors.Is(err, user_service.ErrEmptyEmail) ||
			errors.Is(err, user_service.ErrEmptyPassword) ||
			errors.Is(err, user_service.ErrEmptyRoleCode) ||
			errors.Is(err, user_service.ErrRoleNotFound) ||
			errors.Is(err, user_service.ErrPasswordHashingFailed) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, user_service.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		// Any other error is internal server error
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register user")
	}
	return ctx.JSON(http.StatusCreated, tokens)
}
