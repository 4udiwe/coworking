package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/cowoking/scheduler-service/config"
	consumer_booking "github.com/4udiwe/cowoking/scheduler-service/internal/consumer/booking"
	"github.com/4udiwe/cowoking/scheduler-service/internal/database"
	auth_session_cleaner_producer "github.com/4udiwe/cowoking/scheduler-service/internal/auth_session_cleaner_producer"
	outbox_repository "github.com/4udiwe/cowoking/scheduler-service/internal/repository/outbox"
	timer_repository "github.com/4udiwe/cowoking/scheduler-service/internal/repository/timer"
	scheduler_service "github.com/4udiwe/cowoking/scheduler-service/internal/service/scheduler"
	"github.com/4udiwe/cowoking/scheduler-service/internal/worker"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	postgres *postgres.Postgres

	// Echo
	echoHandler *echo.Echo

	// Repositories
	timerRepo  *timer_repository.TimerRepository
	outboxRepo *outbox_repository.Repository

	// Services
	schedulerService *scheduler_service.SchedulerService

	// Consumer
	bookingConsumer *consumer_booking.Consumer

	// Outbox
	outboxWorker *outbox.Worker

	// Scheduler worker
	scheduerWorker *worker.Worker

	// Session Cleanup Worker
	sessionCleanupWorker *auth_session_cleaner_producer.SessionCleanupWorker
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	return &App{
		interrupt: interrupt,
		cfg:       cfg,
	}
}

func (app *App) Start() {
	// Postgres
	log.Info("Connecting to PostgreSQL...")

	postgres, err := postgres.New(app.cfg.Postgres.URL, postgres.ConnAttempts(5))

	if err != nil {
		log.Fatalf("app - Start - Postgres failed:%v", err)
	}
	app.postgres = postgres

	defer postgres.Close()

	// Migrations
	if err := database.RunMigrations(context.Background(), app.postgres.Pool); err != nil {
		log.Errorf("app - Start - Migrations failed: %v", err)
	}

	// Общий контекст для всех фоновых задач (HTTP‑сервер, Kafka‑consumer’ы, OutboxWorker).
	// По сигналу/ошибке сервера мы вызываем cancel() и корректно останавливаем фоновые горутины.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// App server
	log.Info("Starting app...")

	// Run consumers and workers
	app.BookingConsumer().Run(ctx)
	app.OutboxWorker().Run(ctx)
	app.ScheduerWorker().Run(ctx)
	app.SessionCleanupWorker().Run(ctx)

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	}

	cancel()

	log.Info("Shutting down...")
}
