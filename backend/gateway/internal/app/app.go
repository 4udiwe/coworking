package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/4udiwe/coworking/gateway/internal/config"
	"github.com/4udiwe/coworking/gateway/internal/router"
)

type App struct {
	server *http.Server
}

func New(cfg *config.Config) *App {

	r := router.New(cfg)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &App{
		server: server,
	}
}

func (a *App) Start() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}