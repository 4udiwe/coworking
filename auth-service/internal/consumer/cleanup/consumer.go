package cleanup

import (
	"context"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/coworking/auth-service/internal/consumer"
	auth_service "github.com/4udiwe/coworking/auth-service/internal/service/auth"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	service  *auth_service.Service
	consumer *kafka.KafkaConsumer
	topic    string
	groupID  string
}

func New(
	service *auth_service.Service,
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
	logrus.Infof(
		"SessionCleanupConsumer: subscribing to topic=%s group=%s",
		c.topic,
		c.groupID,
	)

	return c.consumer.Subscribe(
		ctx,
		c.topic,
		c.groupID,
		func(ctx context.Context, key, value []byte) error {
			event, err := consumer.ParseCleanupEvent(value)
			if err != nil {
				logrus.WithError(err).Error(
					"SessionCleanupConsumer: failed to parse event",
				)
				return nil // Skip malformed messages
			}

			if event.Type != consumer.SessionsCleanup {
				logrus.WithFields(logrus.Fields{
					"event_type":  event.Type,
				}).Warn("SessionCleanupConsumer: ignoring unknown event type")
				return nil
			}

			return c.cleanupSessions(ctx, event.Payload.RetentionDays)
		},
	)
}

func (c *Consumer) cleanupSessions(ctx context.Context, retentionDays int) error {
	err := c.service.CleanupOldSessions(ctx, retentionDays)
	if err != nil {
		logrus.WithError(err).WithField("retention_days", retentionDays).Error(
			"SessionCleanupConsumer: failed to cleanup old sessions",
		)
		return err
	}

	return nil
}
