package app

import (
	"github.com/4udiwe/coworking/auth-service/internal/hasher"
	auth_service "github.com/4udiwe/coworking/auth-service/internal/service/auth"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service/user"
)

func (app *App) Hasher() *hasher.BcryptHasher {
	if app.hasher != nil {
		return app.hasher
	}
	app.hasher = hasher.NewBcryptHasher(app.cfg.Hasher.Cost)
	return app.hasher
}

func (app *App) AuthService() *auth_service.Service {
	if app.authService != nil {
		return app.authService
	}
	app.authService = auth_service.New(
		app.UserRepo(),
		app.AuthRepo(),
		app.Postgres(),
		app.Auth(),
		app.Hasher(),
		app.cfg.Auth.RefreshTokenTTL,
	)
	return app.authService
}

func (app *App) UserService() *user_service.Service {
	if app.userService != nil {
		return app.userService
	}
	app.userService = user_service.New(
		app.UserRepo(),
		app.Postgres(),
	)
	return app.userService
}
