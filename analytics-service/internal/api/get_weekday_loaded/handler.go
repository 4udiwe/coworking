package get_weekday_loaded

import (
	"net/http"

	"github.com/4udiwe/coworking/analytics-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s AnalyticsService
}

func New(analyticsService AnalyticsService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: analyticsService})
}

type Request struct {
	CoworkingID uuid.UUID `json:"coworkingId"`
}

type Response struct {
	Load map[int]int `json:"load"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {

	loadMap, err := h.s.GetWeekdayLoad(ctx.Request().Context(), in.CoworkingID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{Load: loadMap})
}
