package app

import (
	log "github.com/sirupsen/logrus"
)

func initLogger(level string) {
	logrusLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(logrusLevel)
	}

	log.SetFormatter(&log.JSONFormatter{
		DisableTimestamp: true,
	})
}
