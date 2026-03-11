package app

import (
	"github.com/4udiwe/coworking/notification-service/internal/api"
	"github.com/4udiwe/coworking/notification-service/internal/api/get_notifications"
	"github.com/4udiwe/coworking/notification-service/internal/api/patch_notification"
	"github.com/4udiwe/coworking/notification-service/internal/api/post_device"
)

func (app *App) GetNotificationsHandler() api.Handler {
	if app.getNotificationsHandler != nil {
		return app.getNotificationsHandler
	}
	app.getNotificationsHandler = get_notifications.New(app.NotificationService())
	return app.getNotificationsHandler
}

func (app *App) PatchNotificationHandler() api.Handler {
	if app.patchNotificationHandler != nil {
		return app.patchNotificationHandler
	}
	app.patchNotificationHandler = patch_notification.New(app.NotificationService())
	return app.patchNotificationHandler
}

func (app *App) PostDeviceHandler() api.Handler {
	if app.postDeviceHandler != nil {
		return app.postDeviceHandler
	}
	app.postDeviceHandler = post_device.New(app.NotificationService())
	return app.postDeviceHandler
}
