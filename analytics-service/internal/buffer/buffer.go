package batch_buffer

import (
	"context"
	"sync"
	"time"

	"github.com/4udiwe/coworking/analytics-service/internal/entity"
	analytics_service "github.com/4udiwe/coworking/analytics-service/internal/service/analytics"
	"github.com/sirupsen/logrus"
)

type BatchBuffer struct {
	mu sync.Mutex

	events []entity.BookingEvent

	batchSize     int
	flushInterval time.Duration

	analytics *analytics_service.AnalyticsService

	ctx context.Context
}

func NewBatchBuffer(
	ctx context.Context,
	batchSize int,
	flushInterval time.Duration,
	analytics *analytics_service.AnalyticsService,
) *BatchBuffer {

	b := &BatchBuffer{
		ctx:           ctx,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		analytics:     analytics,
	}

	go b.flushLoop()

	return b
}

func (b *BatchBuffer) Add(event entity.BookingEvent) {
	b.mu.Lock()

	b.events = append(b.events, event)
	shouldFlush := len(b.events) >= b.batchSize

	b.mu.Unlock()

	if shouldFlush {
		b.flush()
	}
}

func (b *BatchBuffer) flushLoop() {

	ticker := time.NewTicker(b.flushInterval)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			b.flush()

		case <-b.ctx.Done():
			logrus.Info("BatchBuffer stopping, flushing remaining events")
			b.flush()
			return
		}
	}
}

func (b *BatchBuffer) flush() {
	logrus.Debugf("BatchBuffer - flushing events amount: %v", len(b.events))

	b.mu.Lock()

	if len(b.events) == 0 {
		b.mu.Unlock()
		return
	}

	events := b.events
	b.events = nil

	b.mu.Unlock()

	ctx, cancel := context.WithTimeout(b.ctx, 5*time.Second)
	defer cancel()

	if err := b.analytics.InsertEvents(ctx, events); err != nil {
		logrus.WithError(err).Errorf("batch insert failed, events lost: %d", len(events))
	}
}
