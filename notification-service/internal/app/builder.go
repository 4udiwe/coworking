package app

import notification_builder "github.com/4udiwe/coworking/notification-service/internal/builder"

func (app *App) NotificationBuilder() *notification_builder.DefaultBuilder {
	if app.notificationBuilder != nil {
		return app.notificationBuilder
	}
	app.notificationBuilder = notification_builder.New()
	return app.notificationBuilder
}
