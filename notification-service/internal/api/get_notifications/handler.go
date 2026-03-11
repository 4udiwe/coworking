package get_notifications

import (
	"net/http"

	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/4udiwe/coworking/notification-service/internal/api"
	"github.com/4udiwe/coworking/notification-service/internal/api/dto"
	"github.com/4udiwe/coworking/notification-service/internal/api/middleware"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const NOTIFICATIONS_LIMIT = 10

type handler struct {
	s NotificationService
}

func New(notificationService NotificationService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: notificationService})
}

type Request = dto.MarkNotificationReadRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {

	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	notifications, err := h.s.FetchUnreadNotifications(ctx.Request().Context(), claims.UserID, NOTIFICATIONS_LIMIT)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, dto.NotificationsResponse{
		Notifications: lo.Map(notifications, func(n entity.Notification, _ int) dto.Notification {
			return dto.Notification{
				ID:        n.ID,
				UserID:    n.UserID,
				Payload:   n.Payload,
				Type:      string(n.Type),
				Title:     n.Title,
				Body:      n.Body,
				IsRead:    n.IsRead,
				CreatedAt: n.CreatedAt,
				ReadAt:    n.ReadAt,
			}
		}),
	})
}
