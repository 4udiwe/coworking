package app

import (
	"fmt"
	"net/http"

	"github.com/4udiwe/coworking/auth-service/internal/api/middleware"
	"github.com/4udiwe/subscription-service/pkg/validator"
	"github.com/labstack/echo/v4"
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

	authGroup := handler.Group("auth")
	{
		authGroup.POST("/login", app.PostLoginHandler().Handle)
		authGroup.POST("/logout", app.PostLogoutHandler().Handle, app.AuthMiddleware().Middleware)
		authGroup.POST("/refresh", app.PostRefreshHandler().Handle)
		authGroup.POST("/register", app.PostRegisterHandler().Handle)
	}

	userGroup := handler.Group("users", app.AuthMiddleware().Middleware)
	{
		userGroup.GET("/me", app.GetMeHandler().Handle)
		userGroup.GET("/sessions/active", app.GetActiveSessionsHandler().Handle)
		userGroup.GET("/sessions/all", app.GetAllSessionsHandler().Handle)
		userGroup.POST("/sessions/revoke", app.PostRevokeSessionHandler().Handle)
	}

	adminGroup := handler.Group("admin", app.AuthMiddleware().Middleware, middleware.AdminOnly)
	{
		adminGroup.GET("/users", app.GetUsersHandler().Handle)
		adminGroup.GET("/users/:userId", app.GetUserByIdHandler().Handle)
		adminGroup.PATCH("/users/:userId/set_active", app.PatchUserSetActiveHandler().Handle)
		adminGroup.PUT("/users/:userId/roles", app.PutUserRolesHandler().Handle)
	}

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
}
