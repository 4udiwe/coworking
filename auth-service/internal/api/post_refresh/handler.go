package post_refresh

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
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	userAgent := ctx.Request().UserAgent()
	ip := ctx.RealIP()

	deviceName := api.ExtractDeviceName(userAgent)

	tokens, err := h.s.Refresh(ctx.Request().Context(), in.RefreshToken, userAgent, deviceName, ip)

	if err != nil {
		// Validation errors
		if errors.Is(err, user_service.ErrEmptyToken) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		// Authentication errors
		if errors.Is(err, user_service.ErrInvalidRefreshToken) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		// Any other error is internal server error
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, tokens)
}
