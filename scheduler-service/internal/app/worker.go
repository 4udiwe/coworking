package app

import (
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/outbox"
	"github.com/4udiwe/cowoking/scheduler-service/internal/worker"
)

func (app *App) ScheduerWorker() *worker.Worker {
	if app.scheduerWorker != nil {
		return app.scheduerWorker
	}
	app.scheduerWorker = worker.NewWorker(
		app.TimerRepo(),
		app.OutboxRepo(),
		app.Postgres(),
		app.cfg.Worker.WorkerBatchLimit,
		app.cfg.Worker.WorkerInterval,
	)
	return app.scheduerWorker
}

func (app *App) OutboxWorker() *outbox.Worker {
	if app.outboxWorker != nil {
		return app.outboxWorker
	}
	kafkaPublisher := kafka.NewKafkaPublisher(app.cfg.Kafka.Brokers)
	app.outboxWorker = outbox.NewWorker(
		app.OutboxRepo(),
		kafkaPublisher,
		app.cfg.Outbox.Topic,
		app.cfg.Outbox.BatchLimit,
		app.cfg.Outbox.RequeBatchLimit,
		app.cfg.Outbox.Interval,
		app.cfg.Outbox.RequeInterval,
	)
	return app.outboxWorker
}
