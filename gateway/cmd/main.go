package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/4udiwe/coworking/gateway/internal/app"
	"github.com/4udiwe/coworking/gateway/internal/config"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Info("Loading config...")
	cfg, err := config.Load("/config/config.yaml")
	if err != nil {
		log.WithError(err).Fatal("Load config fail")
	}

	a := app.New(cfg)

	go func() {
		log.Info("Starting app...")
		if err := a.Start(); err != nil {
			log.WithError(err).Fatal("App start fail")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	
	log.Info("Shutting down...")
	a.Shutdown(context.Background())
}
