package worker

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// MediaService определяет интерфейс для работы со старыми медиа, которые зависли в статусе Processing.
type MediaService interface {
	HandleStale(ctx context.Context, limit int) error
}

// StaleChecker периодически проверяет медиа в статусе Processing, которые висят там слишком долго (например, больше 30 минут).
// Если таких медиа много, то может обрабатывать их пачками (например, по 50 штук за раз).
// Для каждого такого медиа:
//
// 1. Если retry_count >= MaxRetryCount → помечаем статусом Failed, чтобы не пытаться обрабатывать бесконечно.
//
// 2. Иначе увеличиваем retry_count на единицу и перезапускаем обработку (например, отправляем в очередь на ресайз).
//
// Это позволяет автоматически "отлавливать" медиа, которые по каким-то причинам зависли в статусе Processing,
// и либо повторять попытки обработки, либо помечать их как Failed после определённого количества неудачных попыток.
type StaleChecker struct {
	service  MediaService
	interval time.Duration
	limit    int

	cancel context.CancelFunc
	done   chan struct{}
}

func NewStaleChecker(
	usecase MediaService,
	interval time.Duration,
	limit int,
) *StaleChecker {
	return &StaleChecker{
		service:  usecase,
		interval: interval,
		limit:    limit,
	}
}

// Start запускает StaleChecker в отдельной горутине. Он будет работать до тех пор, пока не будет вызван Stop или не будет отменён контекст.
// Если Start уже был вызван, то он не запустит новый цикл и выведет предупреждение в лог.
func (w *StaleChecker) Start(parentCtx context.Context) {

	if w.done != nil {
		logrus.Warn("stale checker already started")
		return
	}

	ctx, cancel := context.WithCancel(parentCtx)

	w.cancel = cancel
	w.done = make(chan struct{})

	go w.run(ctx)
}

// run запускает бесконечный цикл, который каждые s.interval вызывает HandleStale для обработки stale media.
func (s *StaleChecker) run(ctx context.Context) {

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	logrus.Info("stale checker started")

	for {
		select {
		case <-ctx.Done():
			logrus.Info("stale checker stopped")
			return

		case <-ticker.C:
			s.handle(ctx)
		}
	}
}

// handle выполняет одну итерацию проверки stale media.
func (w *StaleChecker) handle(ctx context.Context) {
	// ограничим время одной итерации
	iterCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := w.service.HandleStale(iterCtx, w.limit); err != nil {
		logrus.WithError(err).Error("stale checker iteration failed")
	}
}

// Stop останавливает StaleChecker, отменяя контекст и ожидая завершения текущей итерации.
func (w *StaleChecker) Stop() {

	if w.cancel == nil {
		return
	}

	w.cancel()

	<-w.done // ждём завершения

	logrus.Info("stale checker stopped")
}
