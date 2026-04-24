package consumer_scheduler

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/cowoking/booking-service/internal/consumer"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
)

// Обработчик событий для топика scheduler
type Consumer struct {
	service  *booking_service.BookingService
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	service *booking_service.BookingService,
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
	logrus.Infof("SchedulerConsumer: subscribing to topic=%s group=%s", c.topic, c.groupID)

	return c.consumer.Subscribe(ctx, c.topic, c.groupID, func(ctx context.Context, key, value []byte) error {
		event, err := consumer.ParseOrderEvent(value)
		if err != nil {
			logrus.Errorf("SchedulerConsumer: failed to parse event: %v", err)
			return nil
		}

		switch event.Type {

		case consumer.BookingExpire:
			err = c.service.CompleteBooking(ctx, event.Payload.BookingID)
			if err != nil {
				logrus.Errorf("SchedulerConsumer: CompleteBooking failed: %v", err)
			}

		default:
			logrus.Errorf("SchedulerConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
