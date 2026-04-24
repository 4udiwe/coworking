package app

import booking_service "github.com/4udiwe/cowoking/booking-service/internal/service/booking"

func (app *App) BookingService() *booking_service.BookingService {
	if app.bookingService != nil {
		return app.bookingService
	}
	app.bookingService = booking_service.New(
		app.BookingRepo(),
		app.PlaceRepo(),
		app.CoworkingRepo(),
		app.OutboxRepo(),
		*app.LayoutValidator(),
		app.Postgres(),
	)
	return app.bookingService
}
