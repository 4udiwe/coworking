package consumer_booking

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/cowoking/scheduler-service/internal/consumer"
	scheduler_service "github.com/4udiwe/cowoking/scheduler-service/internal/service/scheduler"
)

// Обработчик событий для топика scheduler
type Consumer struct {
	service  *scheduler_service.SchedulerService
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	service *scheduler_service.SchedulerService,
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
			logrus.Errorf("BookingConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {

		case consumer.BookingCreated:
			err = c.service.HandleCreatedBooking(
				ctx,
				event.Payload.BookingID,
				event.Payload.UserID,
				event.Payload.PlaceID,
				event.Payload.StartTime,
				event.Payload.EndTime,
			)
			if err != nil {
				logrus.Errorf("BookingConsumer: HandleCreatedBooking failed: %v", err)
			}

		case consumer.BookingCancelled:
			err = c.service.HandleCancelledBooking(
				ctx,
				event.Payload.BookingID,
			)
			if err != nil {
				logrus.Errorf("BookingConsumer: HandleCancelledBooking failed: %v", err)
			}

		default:
			logrus.Errorf("BookingConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
