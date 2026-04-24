package app

import (
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	device_repository "github.com/4udiwe/coworking/notification-service/internal/repository/device"
	notification_repository "github.com/4udiwe/coworking/notification-service/internal/repository/notification"
	outbox_repository "github.com/4udiwe/coworking/notification-service/internal/repository/outbox"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) DeviceRepo() *device_repository.DeviceRepository {
	if app.deviceRepo != nil {
		return app.deviceRepo
	}
	app.deviceRepo = device_repository.New(app.Postgres())
	return app.deviceRepo
}

func (app *App) NotificationRepo() *notification_repository.NotificationRepository {
	if app.notificationRepo != nil {
		return app.notificationRepo
	}
	app.notificationRepo = notification_repository.New(app.Postgres())
	return app.notificationRepo
}

func (app *App) OutboxRepo() *outbox_repository.Repository {
	if app.outboxRepo != nil {
		return app.outboxRepo
	}
	app.outboxRepo = outbox_repository.New(app.Postgres())
	return app.outboxRepo
}
