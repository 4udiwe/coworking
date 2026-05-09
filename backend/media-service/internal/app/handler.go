package app

import (
	api "github.com/4udiwe/coworking/backend/media-service/internal/api/http"
	"github.com/4udiwe/coworking/backend/media-service/internal/api/http/delete_media"
	"github.com/4udiwe/coworking/backend/media-service/internal/api/http/post_media"
)

func (app *App) PostMediaHandler() api.Handler {
	if app.postMediaHandler != nil {
		return app.postMediaHandler
	}
	app.postMediaHandler = post_media.New(app.mediaService)
	return app.postMediaHandler
}

func (app *App) DeleteMediaHandler() api.Handler {
	if app.deleteMediaHandler != nil {
		return app.deleteMediaHandler
	}
	app.deleteMediaHandler = delete_media.New(app.mediaService)
	return app.deleteMediaHandler
}
