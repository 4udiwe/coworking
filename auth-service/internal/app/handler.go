package app

import (
	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_login"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_logout"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_refresh"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_register"
)

func (app *App) PostLoginHandler() api.Handler {
	if app.postLoginHandler != nil {
		return app.postLoginHandler
	}
	app.postLoginHandler = post_login.New(app.UserService())
	return app.postLoginHandler
}

func (app *App) PostLogoutHandler() api.Handler {
	if app.postLogoutHandler != nil {
		return app.postLogoutHandler
	}
	app.postLogoutHandler = post_logout.New(app.UserService())
	return app.postLogoutHandler
}

func (app *App) PostRefreshHandler() api.Handler {
	if app.postRefreshHandler != nil {
		return app.postRefreshHandler
	}
	app.postRefreshHandler = post_refresh.New(app.UserService())
	return app.postRefreshHandler
}

func (app *App) PostRegisterHandler() api.Handler {
	if app.postRegisterHandler != nil {
		return app.postRegisterHandler
	}
	app.postRegisterHandler = post_register.New(app.UserService())
	return app.postRegisterHandler
}
