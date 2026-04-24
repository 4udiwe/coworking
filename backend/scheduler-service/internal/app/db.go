package app

import (
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	outbox_repository "github.com/4udiwe/cowoking/scheduler-service/internal/repository/outbox"
	timer_repository "github.com/4udiwe/cowoking/scheduler-service/internal/repository/timer"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) TimerRepo() *timer_repository.TimerRepository {
	if app.timerRepo != nil {
		return app.timerRepo
	}
	app.timerRepo = timer_repository.New(app.Postgres())
	return app.timerRepo
}

func (app *App) OutboxRepo() *outbox_repository.Repository {
	if app.outboxRepo != nil {
		return app.outboxRepo
	}
	app.outboxRepo = outbox_repository.New(app.Postgres())
	return app.outboxRepo
}
