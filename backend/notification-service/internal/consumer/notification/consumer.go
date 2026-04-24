package consumer_notification

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/coworking/notification-service/internal/consumer"
	notification_service "github.com/4udiwe/coworking/notification-service/internal/service/notification"
)

// Обработчик событий для топика notification
type Consumer struct {
	service  *notification_service.NotificationService
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	service *notification_service.NotificationService,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		service:  service,
		consumer: consumer,
		topic:    topic,
		groupID:  groupID,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	logrus.Infof("BookingConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		event, err := consumer.ParseOrderEvent(value)
		if err != nil {
			logrus.Errorf("NotificationConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {

		case consumer.NotificationCreated:
			err := c.service.NotifyUser(ctx, event.Payload.NotificationID)
			if err != nil {
				logrus.Errorf("NotificationConsumer: NotificationCreated.NotifyUser failed: %v", err)
			}

		default:
			logrus.Errorf("NotificationConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
