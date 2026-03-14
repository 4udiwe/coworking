package app

import (
	analytics_repository "github.com/4udiwe/coworking/analytics-service/internal/repository/analytics"
	"github.com/4udiwe/coworking/analytics-service/pkg/clickhouse"
)

func (app *App) ClickHouse() *clickhouse.ClickHouse {
	return app.clickhouse
}

func (app *App) AnalyticsRepo() *analytics_repository.AnalyticsRepository {
	if app.analyticsRepo != nil {
		return app.analyticsRepo
	}
	app.analyticsRepo = analytics_repository.New(app.ClickHouse())
	return app.analyticsRepo
}
