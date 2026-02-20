package get_active_sessions

import (
	"errors"
	"net/http"

	api "github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/middleware"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service"
	"github.com/4udiwe/coworking/auth-service/pgk/decorator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const ONLY_ACTIVE = true

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request struct{}

type ResponseSession struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	UserAgent  string `json:"userAgent"`
	Device     string `json:"device,omitempty"`
	IPAddress  string `json:"ipAddress"`
	Revoked    bool   `json:"revoked"`
	CreatedAt  string `json:"createdAt"`
	ExpiresAt  string `json:"expiresAt"`
	LastUsedAt string `json:"lastUsedAt"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	claims, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	sessions, err := h.s.GetUserSessions(ctx.Request().Context(), claims.UserID, ONLY_ACTIVE)

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
	return ctx.JSON(http.StatusOK, lo.Map(sessions, func(session entity.Session, _ int) ResponseSession {
		return ResponseSession{
			ID:         session.ID.String(),
			UserID:     session.UserID.String(),
			UserAgent:  session.UserAgent,
			Device:     lo.FromPtr(session.DeviceName),
			IPAddress:  session.IPAddress,
			Revoked:    session.Revoked,
			CreatedAt:  session.CreatedAt.Format("2006-01-02T15:04:05Z"),
			ExpiresAt:  session.ExpiresAt.Format("2006-01-02T15:04:05Z"),
			LastUsedAt: session.LastUsedAt.Format("2006-01-02T15:04:05Z"),
		}
	}))
}
