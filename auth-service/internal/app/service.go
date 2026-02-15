package app

import (
	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/internal/hasher"
	user_service "github.com/4udiwe/coworking/auth-service/internal/service"
)

func (app *App) Auth() *auth.Auth {
	if app.auth != nil {
		return app.auth
	}
	app.auth = auth.New(
		app.cfg.Auth.AccessTokenSecret,
		app.cfg.Auth.RefreshTokenSecret,
		app.cfg.Auth.AccessTokenTTL,
		app.cfg.Auth.RefreshTokenTTL,
	)
	return app.auth
}

func (app *App) Hasher() *hasher.BcryptHasher {
	if app.hasher != nil {
		return app.hasher
	}
	app.hasher = hasher.NewBcryptHasher(app.cfg.Hasher.Cost)
	return app.hasher
}

func (app *App) UserService() *user_service.Service {
	if app.userService != nil {
		return app.userService
	}
	app.userService = user_service.New(
		app.UserRepo(),
		app.AuthRepo(),
		app.Postgres(),
		app.Auth(),
		app.Hasher(),
	)
	return app.userService
}
