package consumer_booking

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	batch_buffer "github.com/4udiwe/coworking/analytics-service/internal/buffer"
	"github.com/4udiwe/coworking/analytics-service/internal/consumer"
	"github.com/4udiwe/coworking/analytics-service/internal/entity"
)

// Обработчик событий для топика booking
type Consumer struct {
	buffer   *batch_buffer.BatchBuffer
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	buffer *batch_buffer.BatchBuffer,
	consumer *kafka.KafkaConsumer,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		buffer:   buffer,
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

		bookingEvent := entity.BookingEvent{
			EventID:     uuid.New(),
			EventType:   string(event.Type),
			BookingID:   event.Payload.BookingID,
			CoworkingID: event.Payload.CoworkingID,
			UserID:      event.Payload.UserID,
			StartTime:   event.Payload.StartTime,
			EndTime:     event.Payload.EndTime,
			Occurred:    event.OccurredAt,
		}
		logrus.WithFields(logrus.Fields{
			"type": event.Type,
		}).Debug("handling event")

		switch event.Type {

		case consumer.BookingCreated:
			bookingEvent.BookingStatus = "created"
			c.buffer.Add(bookingEvent)

		case consumer.BookingCancelled:
			bookingEvent.BookingStatus = "cancelled"
			c.buffer.Add(bookingEvent)

		case consumer.BookingCompleted:
			bookingEvent.BookingStatus = "completed"
			c.buffer.Add(bookingEvent)

		default:
			logrus.Errorf("BookingConsumer: unknown event type %s", event.Type)
			return nil
		}

		return err
	})
}
