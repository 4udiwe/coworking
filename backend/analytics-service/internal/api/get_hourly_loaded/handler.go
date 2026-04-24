package get_hourly_loaded

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/coworking/analytics-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type handler struct {
	s AnalyticsService
}

func New(analyticsService AnalyticsService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: analyticsService})
}

type Request struct {
	CoworkingID uuid.UUID `param:"coworkingId"`
	Weekday     *int      `query:"weekday"`
}

// Validate проверяет корректность параметров запроса
func (r Request) Validate() error {
	logrus.Debug("Validating query")
	if r.Weekday != nil {
		if *r.Weekday < 1 || *r.Weekday > 7 {
			return fmt.Errorf("weekday must be between 1 and 7 (1=Monday, 7=Sunday), got %d", *r.Weekday)
		}
	}
	return nil
}

type Response struct {
	Load map[int]int `json:"load"`
}

func (h *handler) Handle(ctx echo.Context, in Request) error {
	in.Validate()

	loadMap, err := h.s.GetHourlyLoad(ctx.Request().Context(), in.CoworkingID, in.Weekday)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, Response{Load: loadMap})
}
