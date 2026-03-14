package app

import (
	"context"
	"crypto/rsa"
	"os"

	"github.com/4udiwe/avito-pvz/pkg/httpserver"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/coworking/analytics-service/config"
	"github.com/4udiwe/coworking/analytics-service/internal/api"
	"github.com/4udiwe/coworking/analytics-service/internal/api/middleware"
	batch_buffer "github.com/4udiwe/coworking/analytics-service/internal/buffer"
	consumer_booking "github.com/4udiwe/coworking/analytics-service/internal/consumer/booking"
	"github.com/4udiwe/coworking/analytics-service/internal/database"
	analytics_repository "github.com/4udiwe/coworking/analytics-service/internal/repository/analytics"
	analytics_service "github.com/4udiwe/coworking/analytics-service/internal/service/analytics"
	"github.com/4udiwe/coworking/analytics-service/pkg/clickhouse"
	"github.com/4udiwe/coworking/auth-service/pkg/jwt_validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	clickhouse *clickhouse.ClickHouse

	// Echo
	echoHandler *echo.Echo

	// Repositories
	analyticsRepo *analytics_repository.AnalyticsRepository

	// Services
	analyticsService *analytics_service.AnalyticsService

	// Handlers
	getCoworkingHeatmapHander api.Handler
	getPlaceHeatmapHander     api.Handler
	getHourlyLoadedHandler    api.Handler
	getWeekdayLoadedHandler   api.Handler

	// Consumer
	bookingConsumer *consumer_booking.Consumer

	// Batch buffer
	buffer *batch_buffer.BatchBuffer

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
	log.Info("Connecting to Clickhouse...")
	migrationDB, err := clickhouse.NewMigrationDB(
		app.cfg.ClickHouse.Addr,
		app.cfg.ClickHouse.DB,
		app.cfg.ClickHouse.User,
		app.cfg.ClickHouse.Pass,
	)
	if err != nil {
		log.Fatalf("app - Start - Clickhouse migrationDB conntect failed:%v", err)
	}

	// Migrations
	if err := database.RunMigrations(context.Background(), migrationDB); err != nil {
		log.Fatalf("app - Start - Migrations failed: %v", err)
	}

	chDB, err := clickhouse.New(
		app.cfg.ClickHouse.Addr,
		app.cfg.ClickHouse.DB,
		app.cfg.ClickHouse.User,
		app.cfg.ClickHouse.Pass,
	)
	if err != nil {
		log.Fatalf("app - Start - Clickhouse connect failed:%v", err)
	}

	app.clickhouse = chDB
	defer chDB.Close()

	// Общий контекст для всех фоновых задач (HTTP‑сервер, Kafka‑consumer’ы, OutboxWorker).
	// По сигналу/ошибке сервера мы вызываем cancel() и корректно останавливаем фоновые горутины.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Batch buffer
	app.buffer = batch_buffer.NewBatchBuffer(
		ctx, app.cfg.BatchBuffer.BatchSize,
		app.cfg.BatchBuffer.FlushInterval,
		app.AnalyticsService(),
	)

	// Consumers
	bookingKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)

	app.bookingConsumer = consumer_booking.New(
		app.buffer,
		bookingKafkaConsumer,
		app.cfg.Kafka.Topics.BookingEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)

	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	// Run consumers
	app.bookingConsumer.Run(ctx)

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
