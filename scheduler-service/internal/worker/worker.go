package worker

import (
	"context"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/transactor"
	outbox_repository "github.com/4udiwe/cowoking/scheduler-service/internal/repository/outbox"
	timer_repository "github.com/4udiwe/cowoking/scheduler-service/internal/repository/timer"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	timerRepo  *timer_repository.TimerRepository
	outboxRepo *outbox_repository.Repository
	txManager  transactor.Transactor

	batchLimit int
	interval   time.Duration
}

func NewWorker(
	timerRepo *timer_repository.TimerRepository,
	outboxRepo *outbox_repository.Repository,
	txManager transactor.Transactor,
	batchLimit int,
	interval time.Duration,
) *Worker {
	return &Worker{
		timerRepo:  timerRepo,
		outboxRepo: outboxRepo,
		txManager:  txManager,
		batchLimit: batchLimit,
		interval:   interval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	go func() {
		logrus.Infof(
			"SchedulerWorker started interval=%s batchLimit=%d",
			w.interval,
			w.batchLimit,
		)

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logrus.Info("SchedulerWorker stopped")
				return

			case <-ticker.C:
				w.processBatch(ctx)
			}
		}
	}()
}

func (w *Worker) processBatch(ctx context.Context) {

	err := w.txManager.WithinTransaction(ctx, func(ctx context.Context) error {

		// 1️⃣ выбираем готовые таймеры
		timers, err := w.timerRepo.FindDueTimers(ctx, w.batchLimit)
		if err != nil {
			logrus.WithError(err).
				Error("SchedulerWorker: failed to fetch due timers")

			return err
		}

		if len(timers) == 0 {
			return nil
		}

		logrus.WithField("count", len(timers)).
			Info("SchedulerWorker: timers fetched")

		mapper := TimerToEventMapper{}

		var triggeredIDs []uuid.UUID

		// 2️⃣ создаем outbox события
		for _, timer := range timers {

			ev := mapper.Map(timer)

			err := w.outboxRepo.Create(ctx, ev)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"timer_id": timer.ID,
					"type":     ev.EventType,
				}).WithError(err).
					Error("SchedulerWorker: failed to create outbox event")

				return err
			}

			triggeredIDs = append(triggeredIDs, timer.ID)

			logrus.WithFields(logrus.Fields{
				"timer_id": timer.ID,
				"type":     ev.EventType,
			}).Info("SchedulerWorker: outbox event created")
		}

		// 3️⃣ обновляем таймеры
		err = w.timerRepo.MarkTriggered(ctx, triggeredIDs)
		if err != nil {
			logrus.WithError(err).
				Error("SchedulerWorker: failed to mark timers triggered")

			return err
		}

		logrus.WithField("count", len(triggeredIDs)).
			Info("SchedulerWorker: timers marked triggered")

		return nil
	})

	if err != nil {
		logrus.WithError(err).
			Error("SchedulerWorker: batch failed")
	}
}
