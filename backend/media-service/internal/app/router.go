package app

import (
	"net/http"
	"strings"

	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/4udiwe/coworking/backend/media-service/internal/api/http/middleware"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()

	app.configureRouter(handler)

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	// Health check endpoint (no auth required)
	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	//Auth middleware with skipper for /health
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

	mediaGroup := handler.Group("/admin/media", middleware.AdminOnly)
	{
		mediaGroup.POST("/upload", app.PostMediaHandler().Handle)
		mediaGroup.DELETE("/:id", app.DeleteMediaHandler().Handle)
	}
}
