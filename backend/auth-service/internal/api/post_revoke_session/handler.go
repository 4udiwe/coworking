package post_revoke_session

import (
	"net/http"

	api "github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s UserService
}

func New(userService UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: userService})
}

type Request struct {
	SessionID uuid.UUID `json:"sessionId" validate:"required"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	err := h.s.RevokeSession(ctx.Request().Context(), in.SessionID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusAccepted)
}
