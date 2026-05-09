package app

import (
	"context"
	"crypto/rsa"
	"os"

	"github.com/4udiwe/avito-pvz/pkg/httpserver"
	"github.com/4udiwe/coworking/auth-service/pkg/jwt_validator"
	"github.com/4udiwe/coworking/backend/media-service/config"
	api "github.com/4udiwe/coworking/backend/media-service/internal/api/http"
	"github.com/4udiwe/coworking/backend/media-service/internal/api/http/middleware"
	"github.com/4udiwe/coworking/backend/media-service/internal/image_processor"
	media_repository "github.com/4udiwe/coworking/backend/media-service/internal/repository/media"
	object_repository "github.com/4udiwe/coworking/backend/media-service/internal/repository/object"
	media_service "github.com/4udiwe/coworking/backend/media-service/internal/service/media"
	"github.com/4udiwe/coworking/backend/media-service/pkg/minio"
	mongodb "github.com/4udiwe/coworking/backend/media-service/pkg/mongo"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// DB
	mongodb *mongodb.MongoDB
	minio   *minio.Client

	// Repository
	mediaRepo     *media_repository.MediaRepository
	objectStorage *object_repository.Storage

	// Service
	mediaService *media_service.MediaService

	// Handlers
	deleteMediaHandler       api.Handler
	patchReorderMediaHandler api.Handler
	postMediaHandler         api.Handler

	// Echo
	echoHandler *echo.Echo

	// Middleware
	authMW *middleware.AuthMiddleware

	// Auth
	PublicKey    *rsa.PublicKey
	jwtValidator *jwt_validator.Validator
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		logrus.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	return &App{
		cfg: cfg,
	}
}

func (app *App) Start() {
	// ──────────────────────────────────────────
	// MongoDB
	// ──────────────────────────────────────────
	logrus.Info("Connecting to MongoDB...")
	mongodb, err := mongodb.New(app.cfg.MongoDB.URI, app.cfg.MongoDB.DBName)
	if err != nil {
		logrus.Fatalf("app - Start - MongoDB connection failed: %v", err)
	}
	app.mongodb = mongodb
	defer mongodb.Close()

	// ──────────────────────────────────────────
	// MinIO Object Storage
	// ──────────────────────────────────────────
	logrus.Info("Connecting to MinIO...")
	minio, err := minio.New(
		app.cfg.MinIO.Endpoint,
		app.cfg.MinIO.AccessKey,
		app.cfg.MinIO.SecretKey,
		"media",
	)
	if err != nil {
		logrus.Fatalf("app - Start - MinIO connection failed: %v", err)
	}
	app.minio = minio

	// ──────────────────────────────────────────
	// Repositories
	// ──────────────────────────────────────────
	app.mediaRepo = media_repository.New(app.mongodb.Database)
	app.objectStorage = object_repository.New(app.minio)

	// ──────────────────────────────────────────
	// Services
	// ──────────────────────────────────────────
	app.mediaService = media_service.New(
		app.mediaRepo,
		app.objectStorage,
		image_processor.New(),
	)

	// ──────────────────────────────────────────
	// HTTP Server (Echo)
	// ──────────────────────────────────────────
	logrus.Info("Setting up HTTP server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	logrus.Debugf("Server port: %s", app.cfg.HTTP.Port)

	// ──────────────────────────────────────────
	// Wait for shutdown signal or error
	// ──────────────────────────────────────────
	logrus.Info("All services started. Waiting for shutdown signal...")
	s := <-app.interrupt
	logrus.Infof("Received signal: %v", s)

	// ──────────────────────────────────────────
	// Graceful shutdown
	// ──────────────────────────────────────────
	logrus.Info("Shutting down servers...")

	// HTTP shutdown
	ctx, cancel := context.WithTimeout(context.Background(), app.cfg.Shutdown.Timeout)
	defer cancel()

	if err := app.echoHandler.Shutdown(ctx); err != nil {
		logrus.Warnf("HTTP server shutdown error: %v", err)
	}
	logrus.Info("HTTP server stopped")

	logrus.Info("Application shutdown complete")
}
