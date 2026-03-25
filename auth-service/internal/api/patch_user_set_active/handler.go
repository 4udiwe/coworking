package patch_user_set_active

import (
	"net/http"

	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s UserService
}

func New(s UserService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: s})
}

type Request struct {
	UserID uuid.UUID `param:"userId" validate:"required"`
	Active *bool     `json:"active" validate:"required"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	if err := h.s.SetUserActive(ctx.Request().Context(), in.UserID, *in.Active); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}
