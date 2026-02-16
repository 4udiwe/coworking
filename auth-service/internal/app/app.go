package app

import (
	"context"
	"os"

	"github.com/4udiwe/avito-pvz/pkg/httpserver"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/coworking/auth-service/config"
	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/middleware"
	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/database"
	"github.com/4udiwe/coworking/auth-service/internal/hasher"
	auth_repository "github.com/4udiwe/coworking/auth-service/internal/repository/auth"
	user_repository "github.com/4udiwe/coworking/auth-service/internal/repository/user"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service"
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
	authRepo *auth_repository.AuthRepository
	userRepo *user_repository.UserRepository

	// Services
	userService *user_service.Service

	// Handlers
	postLoginHandler         api.Handler
	postLogoutHandler        api.Handler
	postRefreshHandler       api.Handler
	postRegisterHandler      api.Handler
	postRevokeSessionHandler api.Handler

	getMeHandler             api.Handler
	getAllSessionsHandler    api.Handler
	getActiveSessionsHandler api.Handler

	// Auth
	auth *auth.Auth

	// Hasher
	hasher *hasher.BcryptHasher

	// Middleware
	authMW *middleware.AuthMiddleware
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

	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	// Stop HTTP‑server after signal/error.
	if err := httpServer.Shutdown(); err != nil {
		log.Errorf("HTTP server shutdown error: %v", err)
	}

	log.Info("Shutting down...")
}
