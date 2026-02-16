package app

import "github.com/4udiwe/coworking/auth-service/internal/api/middleware"

func (app *App) AuthMiddleware() *middleware.AuthMiddleware {
	if app.authMW != nil {
		return app.authMW
	}
	app.authMW = middleware.New(app.Auth())
	return app.authMW
}
