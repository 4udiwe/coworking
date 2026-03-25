package app

import (
	"github.com/4udiwe/coworking/auth-service/internal/api"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_active_sessions"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_all_sessions"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_me"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_user_by_id"
	"github.com/4udiwe/coworking/auth-service/internal/api/get_users"
	"github.com/4udiwe/coworking/auth-service/internal/api/patch_user_set_active"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_login"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_logout"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_refresh"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_register"
	"github.com/4udiwe/coworking/auth-service/internal/api/post_revoke_session"
	"github.com/4udiwe/coworking/auth-service/internal/api/put_user_roles"
)

func (app *App) PostLoginHandler() api.Handler {
	if app.postLoginHandler != nil {
		return app.postLoginHandler
	}
	app.postLoginHandler = post_login.New(app.AuthService())
	return app.postLoginHandler
}

func (app *App) PostLogoutHandler() api.Handler {
	if app.postLogoutHandler != nil {
		return app.postLogoutHandler
	}
	app.postLogoutHandler = post_logout.New(app.AuthService())
	return app.postLogoutHandler
}

func (app *App) PostRefreshHandler() api.Handler {
	if app.postRefreshHandler != nil {
		return app.postRefreshHandler
	}
	app.postRefreshHandler = post_refresh.New(app.AuthService())
	return app.postRefreshHandler
}

func (app *App) PostRegisterHandler() api.Handler {
	if app.postRegisterHandler != nil {
		return app.postRegisterHandler
	}
	app.postRegisterHandler = post_register.New(app.AuthService())
	return app.postRegisterHandler
}

func (app *App) PostRevokeSessionHandler() api.Handler {
	if app.postRevokeSessionHandler != nil {
		return app.postRevokeSessionHandler
	}
	app.postRevokeSessionHandler = post_revoke_session.New(app.AuthService())
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
	app.getActiveSessionsHandler = get_active_sessions.New(app.AuthService())
	return app.getActiveSessionsHandler
}

func (app *App) GetAllSessionsHandler() api.Handler {
	if app.getAllSessionsHandler != nil {
		return app.getAllSessionsHandler
	}
	app.getAllSessionsHandler = get_all_sessions.New(app.AuthService())
	return app.getAllSessionsHandler
}

func (app *App) GetUserByIdHandler() api.Handler {
	if app.getUserByIdHandler != nil {
		return app.getUserByIdHandler
	}
	app.getUserByIdHandler = get_user_by_id.New(app.UserService())
	return app.getUserByIdHandler
}

func (app *App) GetUsersHandler() api.Handler {
	if app.getUsersHandler != nil {
		return app.getUsersHandler
	}
	app.getUsersHandler = get_users.New(app.UserService())
	return app.getUsersHandler
}

func (app *App) PatchUserSetActiveHandler() api.Handler {
	if app.patchUserSetActiveHandler != nil {
		return app.patchUserSetActiveHandler
	}
	app.patchUserSetActiveHandler = patch_user_set_active.New(app.UserService())
	return app.patchUserSetActiveHandler
}

func (app *App) PutUserRolesHandler() api.Handler {
	if app.putUserRolesHanlder != nil {
		return app.putUserRolesHanlder
	}
	app.putUserRolesHanlder = put_user_roles.New(app.UserService())
	return app.putUserRolesHanlder
}
