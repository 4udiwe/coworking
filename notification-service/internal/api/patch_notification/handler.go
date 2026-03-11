package patch_notification

import (
	"net/http"

	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/4udiwe/coworking/notification-service/internal/api"
	"github.com/4udiwe/coworking/notification-service/internal/api/dto"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s NotificationService
}

func New(notificationService NotificationService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: notificationService})
}

type Request = dto.MarkNotificationReadRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {

	err := h.s.MarkRead(ctx.Request().Context(), in.NotificationID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}
