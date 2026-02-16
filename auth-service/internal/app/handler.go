package app

import (
	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_active_sessions"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_all_sessions"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_me"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_login"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_logout"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_refresh"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_register"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_revoke_session"
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

func (app *App) PostRevokeSessionHandler() api.Handler {
	if app.postRevokeSessionHandler != nil {
		return app.postRevokeSessionHandler
	}
	app.postRevokeSessionHandler = post_revoke_session.New(app.UserService())
	return app.postRevokeSessionHandler
}

func (app *App) GetMeHandler() api.Handler {
	if app.getMeHandler != nil {
		return app.getMeHandler
	}
	app.getMeHandler = get_me.New(app.UserService())
	return app.getMeHandler
}

func (app *App) GetActiveSessionsHandler() api.Handler {
	if app.getActiveSessionsHandler != nil {
		return app.getActiveSessionsHandler
	}
	app.getActiveSessionsHandler = get_active_sessions.New(app.UserService())
	return app.getActiveSessionsHandler
}

func (app *App) GetAllSessionsHandler() api.Handler {
	if app.getAllSessionsHandler != nil {
		return app.getAllSessionsHandler
	}
	app.getAllSessionsHandler = get_all_sessions.New(app.UserService())
	return app.getAllSessionsHandler
}
