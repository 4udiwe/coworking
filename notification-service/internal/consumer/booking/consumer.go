package consumer_booking

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	notification_builder "github.com/4udiwe/coworking/notification-service/internal/builder"
	"github.com/4udiwe/coworking/notification-service/internal/consumer"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	notification_service "github.com/4udiwe/coworking/notification-service/internal/service/notification"
)

// Обработчик событий для топика booking
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
	logrus.Infof("BookingConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		event, err := consumer.ParseOrderEvent(value)
		if err != nil {
			logrus.Errorf("BookingConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {

		case consumer.BookingCreated:
			builderEvent := notification_builder.Event{
				Type:   entity.BookingCreatedNotificationType,
				UserID: event.Payload.UserID,
				Payload: map[string]any{
					"placeId":   event.Payload.PlaceID,
					"startTime": event.Payload.StartTime,
				},
			}
			notification, err := c.builder.Build(builderEvent)
			if err != nil {
				logrus.Errorf("BookingConsumer: BookingCreated.BuildNotification failed: %v", err)
			}

			err = c.service.CreateNotification(ctx, notification)
			if err != nil {
				logrus.Errorf("BookingConsumer: BookingCreated.CreateNotification failed: %v", err)
			}

		case consumer.BookingCancelled:
			builderEvent := notification_builder.Event{
				Type:   entity.BookingCancelledNotificationType,
				UserID: event.Payload.UserID,
				Payload: map[string]any{
					"placeId": event.Payload.PlaceID,
				},
			}
			notification, err := c.builder.Build(builderEvent)
			if err != nil {
				logrus.Errorf("BookingConsumer: BookingCancelled.BuildNotification failed: %v", err)
			}

			err = c.service.CreateNotification(ctx, notification)
			if err != nil {
				logrus.Errorf("BookingConsumer: BookingCancelled.CreateNotification failed: %v", err)
			}

		case consumer.BookingCompleted:
			builderEvent := notification_builder.Event{
				Type:   entity.BookingExpiredNotificationType,
				UserID: event.Payload.UserID,
				Payload: map[string]any{
					"placeId": event.Payload.PlaceID,
					"endTime": event.Payload.EndTime,
				},
			}
			notification, err := c.builder.Build(builderEvent)
			if err != nil {
				logrus.Errorf("BookingConsumer: BookingCompleted.BuildNotification failed: %v", err)
			}

			err = c.service.CreateNotification(ctx, notification)
			if err != nil {
				logrus.Errorf("BookingConsumer: BookingCompleted.CreateNotification failed: %v", err)
			}

		default:
			logrus.Errorf("BookingConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
