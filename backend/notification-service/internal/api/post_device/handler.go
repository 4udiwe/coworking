package post_device

import (
	"net/http"

	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/4udiwe/coworking/notification-service/internal/api"
	"github.com/4udiwe/coworking/notification-service/internal/api/dto"
	"github.com/4udiwe/coworking/notification-service/internal/api/middleware"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s NotificationService
}

func New(notificationService NotificationService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: notificationService})
}

type Request = dto.RegisterDeviceRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = h.s.RegisterDevice(ctx.Request().Context(), claims.UserID, in.DeviceToken, in.Platform)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}
