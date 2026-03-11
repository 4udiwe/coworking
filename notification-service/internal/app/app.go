package app

import (
	"context"
	"crypto/rsa"
	"os"

	"github.com/4udiwe/avito-pvz/pkg/httpserver"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/coworking/auth-service/pkg/jwt_validator"
	"github.com/4udiwe/coworking/notification-service/config"
	"github.com/4udiwe/coworking/notification-service/internal/api"
	"github.com/4udiwe/coworking/notification-service/internal/api/middleware"
	notification_builder "github.com/4udiwe/coworking/notification-service/internal/builder"
	consumer_booking "github.com/4udiwe/coworking/notification-service/internal/consumer/booking"
	consumer_notification "github.com/4udiwe/coworking/notification-service/internal/consumer/notification"
	consumer_scheduler "github.com/4udiwe/coworking/notification-service/internal/consumer/scheduler"
	database "github.com/4udiwe/coworking/notification-service/internal/database/migrations"
	device_repository "github.com/4udiwe/coworking/notification-service/internal/repository/device"
	notification_repository "github.com/4udiwe/coworking/notification-service/internal/repository/notification"
	outbox_repository "github.com/4udiwe/coworking/notification-service/internal/repository/outbox"
	firebase_sender "github.com/4udiwe/coworking/notification-service/internal/sender/firebase"
	notification_service "github.com/4udiwe/coworking/notification-service/internal/service/notification"
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
	deviceRepo       *device_repository.DeviceRepository
	notificationRepo *notification_repository.NotificationRepository
	outboxRepo       *outbox_repository.Repository

	// Services
	notificationService *notification_service.NotificationService

	// Handlers
	getNotificationsHandler  api.Handler
	patchNotificationHandler api.Handler
	postDeviceHandler        api.Handler

	// Consumer
	schedulerConsumer    *consumer_scheduler.Consumer
	bookingConsumer      *consumer_booking.Consumer
	notificationConsumer *consumer_notification.Consumer

	// Push sender
	pushSender *firebase_sender.FirebaseSender

	// Notification builder
	notificationBuilder *notification_builder.DefaultBuilder

	// Outbox
	OutboxWorker *outbox.Worker

	// Middleware
	authMW *middleware.AuthMiddleware

	// Auth
	PublicKey    *rsa.PublicKey
	jwtValidator *jwt_validator.Validator
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	return &App{
		cfg: cfg,
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

	// Push sender
	app.pushSender, err = firebase_sender.New(ctx, app.cfg.PushSender.ServiceAccountPath)
	if err != nil {
		log.Errorf("app - Start - PushSender create failed: %v", err)
	}

	// Consumers
	schedulerKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)
	app.schedulerConsumer = consumer_scheduler.New(
		app.NotificationService(),
		app.NotificationBuilder(),
		schedulerKafkaConsumer,
		app.cfg.Kafka.Topics.SchedulerEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	bookingKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)
	app.bookingConsumer = consumer_booking.New(
		app.NotificationService(),
		app.NotificationBuilder(),
		bookingKafkaConsumer,
		app.cfg.Kafka.Topics.BookingEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	notificationKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)
	app.notificationConsumer = consumer_notification.New(
		app.NotificationService(),
		notificationKafkaConsumer,
		app.cfg.Kafka.Topics.NotificationEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	// Outbox publisher
	kafkaPublisher := kafka.NewKafkaPublisher(app.cfg.Kafka.Brokers)

	app.OutboxWorker = outbox.NewWorker(
		app.OutboxRepo(),
		kafkaPublisher,
		app.cfg.Outbox.Topic,
		app.cfg.Outbox.BatchLimit,
		app.cfg.Outbox.RequeBatchLimit,
		app.cfg.Outbox.Interval,
		app.cfg.Outbox.RequeInterval,
	)

	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	// Run consumers and publisher
	app.schedulerConsumer.Run(ctx)
	app.bookingConsumer.Run(ctx)
	app.notificationConsumer.Run(ctx)
	app.OutboxWorker.Run(ctx)

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	// Останавливаем HTTP‑сервер после получения сигнала/ошибки.
	if err := httpServer.Shutdown(); err != nil {
		log.Errorf("HTTP server shutdown error: %v", err)
	}

	log.Info("Shutting down...")
}
