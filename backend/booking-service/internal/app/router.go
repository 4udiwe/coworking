package app

import (
	"fmt"
	"net/http"
	"strings"

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
	// Health check endpoint (no auth required)
	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	// Auth middleware with skipper for /health
	authMiddleware := app.AuthMiddleware()
	handler.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip auth for /health endpoint
			if strings.HasPrefix(c.Request().URL.Path, "/health") {
				return next(c)
			}
			return authMiddleware.Middleware(next)(c)
		}
	})

	// Public coworking and layout endpoints
	coworkingGroup := handler.Group("/coworkings")
	{
		coworkingGroup.GET("", app.GetCoworkingsHandler().Handle)
		coworkingGroup.GET("/:coworkingId", app.GetCoworkingByIdHandler().Handle)

		coworkingGroup.GET("/:coworkingId/places", app.GetPlacesByCoworkingHandler().Handle)
		coworkingGroup.GET("/:coworkingId/available-places", app.GetAvailablePlacesByCoworkingHandler().Handle)

		coworkingGroup.GET("/:coworkingId/layout", app.GetLayoutHandler().Handle)

	}

	// Public booking endpoints (user-oriented)
	bookingGroup := handler.Group("/bookings")
	{
		bookingGroup.POST("", app.PostBookingHandler().Handle)
		bookingGroup.GET("/:bookingId", app.GetBookingByIdHandler().Handle)
		bookingGroup.GET("/active", app.GetActiveBookingsByUserHandler().Handle)
		bookingGroup.GET("/history", app.GetHistoryBookingsByUserHandler().Handle)
		bookingGroup.DELETE("/:bookingId", app.DeleteBookingHandler().Handle)
	}

	// Admin endpoints
	adminGroup := handler.Group("/admin", middleware.AdminOnly)
	{
		adminCoworkingGroup := adminGroup.Group("/coworkings")
		{

			adminCoworkingGroup.POST("", app.PostCoworkingHandler().Handle)
			adminCoworkingGroup.PUT("/:coworkingId", app.PutCoworkingHandler().Handle)
			adminCoworkingGroup.PATCH("/:coworkingId/set_active", app.PatchCoworkingActiveHandler().Handle)

			adminCoworkingGroup.GET("/:coworkingId/layouts", app.GetLayoutVersionsHandler().Handle)
			adminCoworkingGroup.POST("/:coworkingId/layouts", app.PostLayoutHandler().Handle)
			adminCoworkingGroup.GET("/:coworkingId/layouts/:version", app.GetLayoutByVersionHandler().Handle)
			adminCoworkingGroup.PATCH("/:coworkingId/layouts/:version", app.PatchLayoutSetActiveHandler().Handle)
			adminCoworkingGroup.DELETE("/:coworkingId/layouts/:version", app.DeleteLayoutHandler().Handle)

		}

		adminPlacesGroup := adminGroup.Group("/places")
		{
			adminPlacesGroup.POST("", app.PostPlacesHandler().Handle)
			adminPlacesGroup.PATCH("/:placeId/set_active", app.PatchPlaceActiveHandler().Handle)
		}

		adminBookingsGroup := adminGroup.Group("/bookings")
		{
			adminBookingsGroup.GET("", app.GetActiveAdminBookingsHandler().Handle)
			adminBookingsGroup.DELETE("/:bookingId", app.DeleteBookingHandler().Handle)
		}
	}
}
