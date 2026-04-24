package app

import (
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	booking_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/booking"
	coworking_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/coworking"
	outbox_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/outbox"
	place_repository "github.com/4udiwe/cowoking/booking-service/internal/repository/place"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) BookingRepo() *booking_repository.BookingRepository {
	if app.bookingRepo != nil {
		return app.bookingRepo
	}
	app.bookingRepo = booking_repository.New(app.Postgres())
	return app.bookingRepo
}

func (app *App) CoworkingRepo() *coworking_repository.CoworkingRepository {
	if app.coworkingRepo != nil {
		return app.coworkingRepo
	}
	app.coworkingRepo = coworking_repository.New(app.Postgres())
	return app.coworkingRepo
}

func (app *App) PlaceRepo() *place_repository.PlaceRepository {
	if app.placeRepo != nil {
		return app.placeRepo
	}
	app.placeRepo = place_repository.New(app.Postgres())
	return app.placeRepo
}

func (app *App) OutboxRepo() *outbox_repository.Repository {
	if app.outboxRepo != nil {
		return app.outboxRepo
	}
	app.outboxRepo = outbox_repository.New(app.Postgres())
	return app.outboxRepo
}
