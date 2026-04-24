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

	message := &messaging.Message{
		Token: msg.Token,

		Notification: &messaging.Notification{
			Title: msg.Title,
			Body:  msg.Body,
		},

		Data: msg.Data,

		Android: &messaging.AndroidConfig{
			Priority: "high",
		},

		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
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
