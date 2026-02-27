package app

import (
	"context"
	"crypto/rsa"
	"os"

	"github.com/4udiwe/avito-pvz/pkg/httpserver"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/cowoking/booking-service/config"
	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/middleware"
	consumer_scheduler "github.com/4udiwe/cowoking/booking-service/internal/consumer/scheduler"
	"github.com/4udiwe/cowoking/booking-service/internal/database"
	booking_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/booking"
	coworking_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/coworking"
	outbox_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/outbox"
	place_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/place"
	booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"
	"github.com/4udiwe/cowoking/booking-service/pkg/json_schema_validator"
	"github.com/4udiwe/coworking/auth-service/pkg/jwt_validator"
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
	coworkingRepo *coworking_repository.CoworkingRepository
	bookingRepo   *booking_repository.BookingRepository
	placeRepo     *place_repository.PlaceRepository
	outboxRepo    *outbox_repository.Repository

	// Services
	bookingService *booking_service.BookingService

	// Handlers
	deleteBookingHandler api.Handler

	getBookingByIdHandler                api.Handler
	getBookingsByUserHandler             api.Handler
	getCoworkingByIdHandler              api.Handler
	getCoworkingsHandler                 api.Handler
	getLayoutHandler                     api.Handler
	getLayoutByVersionHandler            api.Handler
	getLayoutVersionsHandler             api.Handler
	getPlacesByCoworkingHandler          api.Handler
	getAvailablePlacesByCoworkingHandler api.Handler

	postBookingHandler        api.Handler
	postCoworkingHandler      api.Handler
	postLayoutHandler         api.Handler
	postPlacesHandler         api.Handler
	postLayoutRollbackHandler api.Handler

	putCoworkingHandler         api.Handler
	putCoworkingActiveHandler   api.Handler
	putCoworkingInactiveHandler api.Handler
	putPlaceActiveHandler       api.Handler

	// Consumer
	schedulerConsumer *consumer_scheduler.Consumer

	// Outbox
	OutboxWorker *outbox.Worker

	// Layout validator
	layoutValidator *json_schema_validator.Validator

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

	// Consumers
	schedulerKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)

	app.schedulerConsumer = consumer_scheduler.New(
		app.BookingService(),
		schedulerKafkaConsumer,
		app.cfg.Kafka.Topics.SchedulerEvents,
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
