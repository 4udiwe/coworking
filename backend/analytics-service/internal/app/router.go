package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	for _, r := range handler.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	handler.Use(app.AuthMiddleware().Middleware)

	// Public coworking and layout endpoints
	coworkingGroup := handler.Group("/analytics")
	{
		coworkingGroup.GET("/coworking_heatmap/:coworkingId", app.GetCoworkingHeatmapHandler().Handle)
		coworkingGroup.GET("/place_heatmap/:placeId", app.GetPlaceHeatmapHandler().Handle)
		coworkingGroup.GET("/hourly/:coworkingId", app.GetHourlyLoadedHandler().Handle)
		coworkingGroup.GET("/weekday/:coworkingId", app.GetWeekdayLoadedHandler().Handle)
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
}
