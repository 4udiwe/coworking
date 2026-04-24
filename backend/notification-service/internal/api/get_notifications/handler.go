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

type Request = dto.GetNotificationsQuery

func (h *handler) Handle(ctx echo.Context, in Request) error {

	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	limit := in.Limit
	if limit == 0 {
		limit = NOTIFICATIONS_LIMIT
	}

	var notifications []entity.Notification

	// If since parameter is provided, use date-based filtering
	if in.Since != nil {
		notifications, err = h.s.FetchNotificationsAfterDate(
			ctx.Request().Context(),
			claims.UserID,
			*in.Since,
			limit,
			in.Offset,
		)
	} else {
		// Otherwise use regular filtering
		notifications, err = h.s.FetchNotifications(
			ctx.Request().Context(),
			claims.UserID,
			limit,
			in.Offset,
			in.IsRead,
		)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, dto.NotificationsResponse{
		Notifications: lo.Map(notifications, func(n entity.Notification, _ int) dto.Notification {
			return dto.Notification{
				ID:        n.ID,
				UserID:    n.UserID,
				Payload:   n.Payload,
				Type:      string(n.Type),
				Title:     n.Title,
				Body:      n.Body,
				ActionURL: n.ActionURL,
				IsRead:    n.IsRead,
				CreatedAt: n.CreatedAt,
				ReadAt:    n.ReadAt,
			}
		}),
	})
}
