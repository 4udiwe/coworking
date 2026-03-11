package app

import notification_service "github.com/4udiwe/coworking/notification-service/internal/service/notification"

func (app *App) NotificationService() *notification_service.NotificationService {
	if app.notificationService != nil {
		return app.notificationService
	}
	app.notificationService = notification_service.New(
		app.NotificationRepo(),
		app.DeviceRepo(),
		app.OutboxRepo(),
		app.PushSender(),
		app.Postgres(),
	)
	return app.notificationService
}
