package get_unread_count

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

type Request struct{}

func (h *handler) Handle(ctx echo.Context, in Request) error {

	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	count, err := h.s.GetUnreadCount(ctx.Request().Context(), claims.UserID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dto.UnreadCountResponse{
		UnreadCount: count,
	})
}
