package post_login

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/coworking/auth-service/internal/api"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
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
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	userAgent := ctx.Request().UserAgent()
	ip := ctx.RealIP()

	deviceName := api.ExtractDeviceName(userAgent)

	tokens, err := h.s.Login(ctx.Request().Context(), in.Email, in.Password, userAgent, deviceName, ip)

	if err != nil {
		// Validation errors
		if errors.Is(err, user_service.ErrEmptyEmail) ||
			errors.Is(err, user_service.ErrEmptyPassword) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		// Authentication errors
		if errors.Is(err, user_service.ErrInvalidCredentials) ||
			errors.Is(err, user_service.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		// Any other error is internal server error
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, tokens)
}
