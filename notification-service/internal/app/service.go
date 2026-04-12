package app

import (
	firebase_sender "github.com/4udiwe/coworking/notification-service/internal/sender/firebase"
	notification_service "github.com/4udiwe/coworking/notification-service/internal/service/notification"
)

func (app *App) NotificationService() *notification_service.NotificationService {
	if app.notificationService != nil {
		return app.notificationService
	}
	app.notificationService = notification_service.New(
		app.NotificationRepo(),
		app.DeviceRepo(),
		app.OutboxRepo(),
		app.PushService(),
		app.Postgres(),
	)
	return app.notificationService
}

func (app *App) PushService() notification_service.PushService {
	return notification_service.NewDefaultPushService(app.DefaultDispatcher())
}

func (app *App) DefaultDispatcher() *firebase_sender.DefaultDispatcher {
	return firebase_sender.NewDefaultDispatcher(app.PushSender(), app.DeviceRepo())
}
