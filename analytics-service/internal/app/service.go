package app

import (
	analytics_service "github.com/4udiwe/coworking/analytics-service/internal/service/analytics"
)

func (app *App) AnalyticsService() *analytics_service.AnalyticsService {
	if app.analyticsService != nil {
		return app.analyticsService
	}
	app.analyticsService = analytics_service.New(app.AnalyticsRepo())
	return app.analyticsService
}
