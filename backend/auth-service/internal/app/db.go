package app

import (
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	auth_repository "github.com/4udiwe/coworking/auth-service/internal/repository/auth"
	user_repository "github.com/4udiwe/coworking/auth-service/internal/repository/user"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) AuthRepo() *auth_repository.AuthRepository {
	if app.authRepo != nil {
		return app.authRepo
	}
	app.authRepo = auth_repository.New(app.Postgres())
	return app.authRepo
}

func (app *App) UserRepo() *user_repository.UserRepository {
	if app.userRepo != nil {
		return app.userRepo
	}
	app.userRepo = user_repository.New(app.Postgres())
	return app.userRepo
}
