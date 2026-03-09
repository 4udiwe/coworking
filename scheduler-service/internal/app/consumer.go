package app

import (
	"github.com/4udiwe/big-bob-pizza/order-service/pkg/kafka"
	consumer_booking "github.com/4udiwe/cowoking/scheduler-service/internal/consumer/booking"
)

func (app *App) BookingConsumer() *consumer_booking.Consumer {
	if app.bookingConsumer != nil {
		return app.bookingConsumer
	}
	bookingKafkaConsumer := kafka.NewConsumer(app.cfg.Kafka.Brokers)

	app.bookingConsumer = consumer_booking.New(
		app.SchedulerService(),
		bookingKafkaConsumer,
		app.cfg.Kafka.Topics.SchedulerEvents,
		app.cfg.Kafka.Consumer.GroupID,
	)
	return app.bookingConsumer
}
