package consumer_scheduler

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	notification_builder "github.com/4udiwe/coworking/notification-service/internal/builder"
	"github.com/4udiwe/coworking/notification-service/internal/consumer"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	notification_service "github.com/4udiwe/coworking/notification-service/internal/service/notification"
)

// Обработчик событий для топика scheduler
type Consumer struct {
	service *notification_service.NotificationService
	builder *notification_builder.DefaultBuilder

	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	service *notification_service.NotificationService,
	builder *notification_builder.DefaultBuilder,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		service:  service,
		builder:  builder,
		consumer: consumer,
		topic:    topic,
		groupID:  groupID,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	logrus.Infof("SchedulerConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		event, err := consumer.ParseOrderEvent(value)
		if err != nil {
			logrus.Errorf("SchedulerConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {

		case consumer.ReminderTriggered:
			builderEvent := notification_builder.Event{
				Type:   entity.BookingReminderNotificationType,
				UserID: event.Payload.UserID,
				Payload: map[string]any{
					"placeId":   event.Payload.PlaceID,
					"startTime": event.Payload.StartTime,
				},
			}
			notification, err := c.builder.Build(builderEvent)
			if err != nil {
				logrus.Errorf("SchedulerConsumer: ReminderTriggered.BuildNotification failed: %v", err)
			}

			err = c.service.CreateNotification(ctx, notification)
			if err != nil {
				logrus.Errorf("SchedulerConsumer: ReminderTriggered.CreateNotification failed: %v", err)
			}

		default:
			logrus.Errorf("SchedulerConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
