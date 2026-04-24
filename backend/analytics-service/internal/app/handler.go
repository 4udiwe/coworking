package app

import (
	"github.com/4udiwe/coworking/analytics-service/internal/api"
	"github.com/4udiwe/coworking/analytics-service/internal/api/get_coworking_heatmap"
	"github.com/4udiwe/coworking/analytics-service/internal/api/get_hourly_loaded"
	"github.com/4udiwe/coworking/analytics-service/internal/api/get_place_heatmap"
	"github.com/4udiwe/coworking/analytics-service/internal/api/get_weekday_loaded"
)

func (app *App) GetCoworkingHeatmapHandler() api.Handler {
	if app.getCoworkingHeatmapHander != nil {
		return app.getCoworkingHeatmapHander
	}
	app.getCoworkingHeatmapHander = get_coworking_heatmap.New(app.AnalyticsService())
	return app.getCoworkingHeatmapHander
}

func (app *App) GetPlaceHeatmapHandler() api.Handler {
	if app.getPlaceHeatmapHander != nil {
		return app.getPlaceHeatmapHander
	}
	app.getPlaceHeatmapHander = get_place_heatmap.New(app.AnalyticsService())
	return app.getPlaceHeatmapHander
}

func (app *App) GetHourlyLoadedHandler() api.Handler {
	if app.getHourlyLoadedHandler != nil {
		return app.getHourlyLoadedHandler
	}
	app.getHourlyLoadedHandler = get_hourly_loaded.New(app.AnalyticsService())
	return app.getHourlyLoadedHandler
}

func (app *App) GetWeekdayLoadedHandler() api.Handler {
	if app.getWeekdayLoadedHandler != nil {
		return app.getWeekdayLoadedHandler
	}
	app.getWeekdayLoadedHandler = get_weekday_loaded.New(app.AnalyticsService())
	return app.getWeekdayLoadedHandler
}