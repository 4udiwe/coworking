package app

import (
	"github.com/4udiwe/cowoking/booking-service/internal/api"
	"github.com/4udiwe/cowoking/booking-service/internal/api/delete_booking"
	"github.com/4udiwe/cowoking/booking-service/internal/api/delete_layout"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_admin_bookings"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_available_places_by_coworking"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_booking_by_id"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_bookings_by_user"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_coworking_by_id"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_coworkings"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_layout"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_layout_by_version"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_layout_versions"
	"github.com/4udiwe/cowoking/booking-service/internal/api/get_places_by_coworking"
	"github.com/4udiwe/cowoking/booking-service/internal/api/patch_coworking_active"
	"github.com/4udiwe/cowoking/booking-service/internal/api/patch_layout_set_active"
	"github.com/4udiwe/cowoking/booking-service/internal/api/patch_place_active"
	"github.com/4udiwe/cowoking/booking-service/internal/api/post_booking"
	"github.com/4udiwe/cowoking/booking-service/internal/api/post_coworking"
	"github.com/4udiwe/cowoking/booking-service/internal/api/post_layout"
	"github.com/4udiwe/cowoking/booking-service/internal/api/post_places"
	"github.com/4udiwe/cowoking/booking-service/internal/api/put_coworking"
)

func (app *App) DeleteBookingHandler() api.Handler {
	if app.deleteBookingHandler != nil {
		return app.deleteBookingHandler
	}
	app.deleteBookingHandler = delete_booking.New(app.BookingService())
	return app.deleteBookingHandler
}

func (app *App) DeleteLayoutHandler() api.Handler {
	if app.deleteLayoutHander != nil {
		return app.deleteLayoutHander
	}
	app.deleteLayoutHander = delete_layout.New(app.BookingService())
	return app.deleteLayoutHander
}

func (app *App) GetBookingByIdHandler() api.Handler {
	if app.getBookingByIdHandler != nil {
		return app.getBookingByIdHandler
	}
	app.getBookingByIdHandler = get_booking_by_id.New(app.BookingService())
	return app.getBookingByIdHandler
}

func (app *App) GetBookingsByUserHandler() api.Handler {
	if app.getBookingsByUserHandler != nil {
		return app.getBookingsByUserHandler
	}
	app.getBookingsByUserHandler = get_bookings_by_user.New(app.BookingService())
	return app.getBookingsByUserHandler
}

func (app *App) GetCoworkingByIdHandler() api.Handler {
	if app.getCoworkingByIdHandler != nil {
		return app.getCoworkingByIdHandler
	}
	app.getCoworkingByIdHandler = get_coworking_by_id.New(app.BookingService())
	return app.getCoworkingByIdHandler
}

func (app *App) GetCoworkingsHandler() api.Handler {
	if app.getCoworkingsHandler != nil {
		return app.getCoworkingsHandler
	}
	app.getCoworkingsHandler = get_coworkings.New(app.BookingService())
	return app.getCoworkingsHandler
}

func (app *App) GetLayoutHandler() api.Handler {
	if app.getLayoutHandler != nil {
		return app.getLayoutHandler
	}
	app.getLayoutHandler = get_layout.New(app.BookingService())
	return app.getLayoutHandler
}

func (app *App) GetLayoutByVersionHandler() api.Handler {
	if app.getLayoutByVersionHandler != nil {
		return app.getLayoutByVersionHandler
	}
	app.getLayoutByVersionHandler = get_layout_by_version.New(app.BookingService())
	return app.getLayoutByVersionHandler
}

func (app *App) GetLayoutVersionsHandler() api.Handler {
	if app.getLayoutVersionsHandler != nil {
		return app.getLayoutVersionsHandler
	}
	app.getLayoutVersionsHandler = get_layout_versions.New(app.BookingService())
	return app.getLayoutVersionsHandler
}

func (app *App) GetPlacesByCoworkingHandler() api.Handler {
	if app.getPlacesByCoworkingHandler != nil {
		return app.getPlacesByCoworkingHandler
	}
	app.getPlacesByCoworkingHandler = get_places_by_coworking.New(app.BookingService())
	return app.getPlacesByCoworkingHandler
}

func (app *App) GetAvailablePlacesByCoworkingHandler() api.Handler {
	if app.getAvailablePlacesByCoworkingHandler != nil {
		return app.getAvailablePlacesByCoworkingHandler
	}
	app.getAvailablePlacesByCoworkingHandler = get_available_places_by_coworking.New(app.BookingService())
	return app.getAvailablePlacesByCoworkingHandler
}

func (app *App) GetActiveAdminBookingsHandler() api.Handler {
	if app.getAdminActiveBookings != nil {
		return app.getAdminActiveBookings
	}
	app.getAdminActiveBookings = get_admin_bookings.New(app.BookingService())
	return app.getAdminActiveBookings
}

func (app *App) PatchCoworkingActiveHandler() api.Handler {
	if app.patchCoworkingActiveHandler != nil {
		return app.patchCoworkingActiveHandler
	}
	app.patchCoworkingActiveHandler = patch_coworking_active.New(app.BookingService())
	return app.patchCoworkingActiveHandler
}

func (app *App) PatchLayoutSetActiveHandler() api.Handler {
	if app.patchLayoutSetActiveHandler != nil {
		return app.patchLayoutSetActiveHandler
	}
	app.patchLayoutSetActiveHandler = patch_layout_set_active.New(app.BookingService())
	return app.patchLayoutSetActiveHandler
}

func (app *App) PatchPlaceActiveHandler() api.Handler {
	if app.patchPlaceActiveHander != nil {
		return app.patchPlaceActiveHander
	}
	app.patchPlaceActiveHander = patch_place_active.New(app.BookingService())
	return app.patchPlaceActiveHander
}

func (app *App) PostBookingHandler() api.Handler {
	if app.postBookingHandler != nil {
		return app.postBookingHandler
	}
	app.postBookingHandler = post_booking.New(app.BookingService())
	return app.postBookingHandler
}

func (app *App) PostCoworkingHandler() api.Handler {
	if app.postCoworkingHandler != nil {
		return app.postCoworkingHandler
	}
	app.postCoworkingHandler = post_coworking.New(app.BookingService())
	return app.postCoworkingHandler
}

func (app *App) PostLayoutHandler() api.Handler {
	if app.postLayoutHandler != nil {
		return app.postLayoutHandler
	}
	app.postLayoutHandler = post_layout.New(app.BookingService())
	return app.postLayoutHandler
}

func (app *App) PostPlacesHandler() api.Handler {
	if app.postPlacesHandler != nil {
		return app.postPlacesHandler
	}
	app.postPlacesHandler = post_places.New(app.BookingService())
	return app.postPlacesHandler
}

func (app *App) PutCoworkingHandler() api.Handler {
	if app.putCoworkingHandler != nil {
		return app.putCoworkingHandler
	}
	app.putCoworkingHandler = put_coworking.New(app.BookingService())
	return app.putCoworkingHandler
}
