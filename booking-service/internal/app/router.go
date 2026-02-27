package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/4udiwe/cowoking/booking-service/internal/api/middleware"
	"github.com/labstack/echo/v4"
	//echoSwagger "github.com/swaggo/echo-swagger"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	for _, r := range handler.Routes() {
		fmt.Printf("%s %s\n", r.Method, r.Path)
	}

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	// Public coworking and layout endpoints
	coworkingGroup := handler.Group("/coworkings")
	{
		coworkingGroup.GET("", app.GetCoworkingsHandler().Handle)
		coworkingGroup.GET("/:coworkingId", app.GetCoworkingByIdHandler().Handle)

		coworkingGroup.GET("/:coworkingId/places", app.GetPlacesByCoworkingHandler().Handle)
		coworkingGroup.GET("/:coworkingId/available-places", app.GetAvailablePlacesByCoworkingHandler().Handle)

		coworkingGroup.GET("/:coworkingId/layout", app.GetLayoutHandler().Handle)
		coworkingGroup.GET("/:coworkingId/layouts", app.GetLayoutVersionsHandler().Handle, middleware.AdminOnly)
		coworkingGroup.GET("/:coworkingId/layouts/:version", app.GetLayoutByVersionHandler().Handle, middleware.AdminOnly)
	}

	// Public booking endpoints (user-oriented)
	bookingGroup := handler.Group("/bookings")
	{
		bookingGroup.POST("", app.PostBookingHandler().Handle)
		bookingGroup.GET("/:bookingId", app.GetBookingByIdHandler().Handle)
		bookingGroup.GET("", app.GetBookingsByUserHandler().Handle)
		bookingGroup.DELETE("/:bookingId", app.DeleteBookingHandler().Handle)
	}

	// Admin endpoints
	adminGroup := handler.Group("/admin", middleware.AdminOnly)
	{
		adminCoworkingGroup := adminGroup.Group("/coworkings")
		{
			adminCoworkingGroup.POST("", app.PostCoworkingHandler().Handle)
			adminCoworkingGroup.PUT("/:coworkingId", app.PutCoworkingHandler().Handle)
			adminCoworkingGroup.PUT("/:coworkingId/activate", app.PutCoworkingActiveHandler().Handle)
			adminCoworkingGroup.PUT("/:coworkingId/deactivate", app.PutCoworkingInactiveHandler().Handle)

			adminCoworkingGroup.POST("/:coworkingId/places", app.PostPlacesHandler().Handle)
			adminCoworkingGroup.POST("/:coworkingId/layouts", app.PostLayoutHandler().Handle)
			adminCoworkingGroup.POST("/:coworkingId/layouts/rollback", app.PostLayoutRollbackHandler().Handle)
		}

		adminPlacesGroup := adminGroup.Group("/places")
		{
			adminPlacesGroup.PUT("/:placeId/active", app.PutPlaceActiveHandler().Handle)
		}

		adminBookingsGroup := adminGroup.Group("/bookings")
		{
			adminBookingsGroup.DELETE("/:bookingId", app.DeleteBookingHandler().Handle)
		}
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
}
