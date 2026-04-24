package auth_session_cleaner_producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/sirupsen/logrus"
)

const SESSION_CLEANUP_EVENT_TYPE = "sessions.cleanup"

type SessionCleanupWorker struct {
	producer      *kafka.KafkaPublisher
	topic         string
	retentionDays int
	interval      time.Duration
}

func New(
	producer *kafka.KafkaPublisher,
	topic string,
	retentionDays int,
	interval time.Duration,
) *SessionCleanupWorker {
	return &SessionCleanupWorker{
		producer:      producer,
		topic:         topic,
		retentionDays: retentionDays,
		interval:      interval,
	}
}

// Run запускает воркер, который периодически отправляет события
func (w *SessionCleanupWorker) Run(ctx context.Context) {
	logrus.Infof(
		"SessionCleanupWorker started: interval=%s, retention=%d days, topic=%s",
		w.interval,
		w.retentionDays,
		w.topic,
	)

	// Отправить событие при старте
	w.sendCleanupEvent(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("SessionCleanupWorker stopped")
			return
		case <-ticker.C:
			w.sendCleanupEvent(ctx)
		}
	}
}

func (w *SessionCleanupWorker) sendCleanupEvent(ctx context.Context) {
	payload := map[string]interface{}{
		"retentionDays": w.retentionDays,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error(
			"SessionCleanupWorker: failed to marshal payload",
		)
		return
	}

	// envelope := kafka.Envelope{
	// 	EventID:    uuid.New(),
	// 	EventType:  "auth.sessions.cleanup",
	// 	OccurredAt: time.Now(),
	// 	Data:       payloadBytes,
	// }

	// envelopeBytes, err := json.Marshal(envelope)
	// if err != nil {
	// 	logrus.WithError(err).Error(
	// 		"SessionCleanupWorker: failed to marshal envelope",
	// 	)
	// 	return
	// }

	err = w.producer.Publish(ctx, w.topic, SESSION_CLEANUP_EVENT_TYPE, payloadBytes)
	if err != nil {
		logrus.WithError(err).Error(
			"SessionCleanupWorker: failed to publish cleanup event",
		)
		return
	}

	logrus.WithFields(logrus.Fields{
		"topic":          w.topic,
		"retention_days": w.retentionDays,
		"event_type":     SESSION_CLEANUP_EVENT_TYPE,
	}).Info("SessionCleanupWorker: cleanup event published")
}
