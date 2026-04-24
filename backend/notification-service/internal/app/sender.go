package app

import firebase_sender "github.com/4udiwe/coworking/notification-service/internal/sender/firebase"

func (app *App) PushSender() *firebase_sender.FirebaseSender {
	return app.pushSender
}
