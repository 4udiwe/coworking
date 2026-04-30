package firebase_sender

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/4udiwe/coworking/notification-service/internal/sender"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type FirebaseSender struct {
	client *messaging.Client
}

func New(ctx context.Context, serviceAccountPath string) (*FirebaseSender, error) {
	opt := option.WithCredentialsFile(serviceAccountPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseSender{
		client: client,
	}, nil
}
func (s *FirebaseSender) Send(
	ctx context.Context,
	msg sender.PushMessage,
) error {
	logrus.Debug("firebase: sending notification")

	// Формирование data со всеми данными для сообщения без блока Notification (только Data)
	 data := map[string]string{
        "title":          msg.Title,
        "body":           msg.Body,
        "notificationId": msg.NotificationID,
        "actionUrl":      msg.ActionURL,
    }
    
    // Добавляем остальной payload
    for k, v := range msg.Data {
        data[k] = v
    }

	message := &messaging.Message{
		Token: msg.Token,

		Data: data,

		Android: &messaging.AndroidConfig{
			Priority: "high",
		},

		APNS: &messaging.APNSConfig{
            Headers: map[string]string{
                "apns-priority": "10",
                "apns-push-type": "alert",
            },
            Payload: &messaging.APNSPayload{
                Aps: &messaging.Aps{
                    ContentAvailable: true,
                },
            },
        },
	}

	_, err := s.client.Send(ctx, message)
	if err != nil {
		if messaging.IsRegistrationTokenNotRegistered(err) {
			logrus.Warn("firebase: device token not registered")
			return sender.ErrInvalidToken
		}

		if messaging.IsInvalidArgument(err) {
			logrus.Warn("firebase: invalid token format")
			return sender.ErrInvalidToken
		}

		logrus.WithError(err).Error("firebase send failed")

		return err
	}

	return nil
}
