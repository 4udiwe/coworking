package delete_media

import (
	"net/http"

	"github.com/4udiwe/coworking/auth-service/pkg/decorator"
	api "github.com/4udiwe/coworking/backend/media-service/internal/api/http"
	"github.com/4udiwe/coworking/backend/media-service/internal/api/http/dto"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	s MediaService
}

func New(mediaService MediaService) api.Handler {
	return decorator.NewBindAndValidateDerocator(&handler{s: mediaService})
}

type Request = dto.DeleteMediaRequest

func (h *handler) Handle(ctx echo.Context, in Request) error {
	objectID, err := primitive.ObjectIDFromHex(in.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid media ID")
	}
	err = h.s.Delete(ctx.Request().Context(), objectID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusAccepted)
}
