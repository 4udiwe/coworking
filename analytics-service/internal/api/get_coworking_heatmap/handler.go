package get_coworking_heatmap

import (
	"net/http"

	"github.com/4udiwe/coworking/analytics-service/internal/api"
	"github.com/4udiwe/coworking/analytics-service/internal/entity"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	s AnalyticsService
}

func New(analyticsService AnalyticsService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: analyticsService})
}

type Request struct {
	CoworkingID uuid.UUID `param:"coworkingId"`
}

type Cell struct {
	Weekday uint8 `json:"weekday"` // День недели (1-7, понедельник=0)
	Hour    uint8 `json:"hour"`    // Час суток (0-23)
	Count   uint64 `json:"count"`   // Количество бронирований в данную ячейку
}

type Response struct {
	HeatMap []Cell `json:"heatmap"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {

	cells, err := h.s.GetCoworkingHeatmap(ctx.Request().Context(), in.CoworkingID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{HeatMap: lo.Map(cells, func(e entity.HeatmapCell, _ int) Cell {
		return Cell{
			Weekday: e.Weekday,
			Hour:    e.Hour,
			Count:   e.Count,
		}
	})})
}
