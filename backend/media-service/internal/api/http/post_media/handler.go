package post_media

import (
    "io"
    "net/http"
    
    api "github.com/4udiwe/coworking/backend/media-service/internal/api/http"
    "github.com/4udiwe/coworking/backend/media-service/internal/api/http/dto"
    "github.com/4udiwe/coworking/backend/media-service/internal/entity"
    "github.com/labstack/echo/v4"
    "github.com/sirupsen/logrus"
)

type handler struct {
    s MediaService
}

func New(mediaService MediaService) api.Handler {
    return &handler{s: mediaService}
}

// Обработчик без декоратора, напрямую работающий с form-data
func (h *handler) Handle(c echo.Context) error {
    logrus.Infof("HTTP %s %s from %s", c.Request().Method, c.Path(), c.Request().RemoteAddr)
    
    // Привязываем form-data вручную
    var req dto.PostMediaRequest
    
    // Bind form values
    req.UploadedBy = c.FormValue("uploaded_by")
    
    // Получаем файл
    file, err := c.FormFile("file")
    if err != nil {
        logrus.Errorf("Failed to get file: %v", err)
        return echo.NewHTTPError(http.StatusBadRequest, "file is required")
    }
    req.File = file
    
    // Ручная валидация
    if err := validateRequest(req); err != nil {
        logrus.Errorf("Validation failed: %v", err)
        return err
    }
    
    // Читаем файл
    src, err := req.File.Open()
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "failed to open file")
    }
    defer src.Close()
    
    data, err := io.ReadAll(src)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "failed to read file")
    }
    
    uploadInput := entity.UploadInput{
        FileName:    req.File.Filename,
        ContentType: req.File.Header.Get("Content-Type"),
        Data:        data,
        UploadedBy:  req.UploadedBy,
    }
    
    uploadResult, err := h.s.Upload(c.Request().Context(), uploadInput)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(http.StatusCreated, dto.PostMediaResponse{
        ID:     uploadResult.ID,
        Status: uploadResult.Status,
        URLs:   uploadResult.URLs,
    })
}

func validateRequest(req dto.PostMediaRequest) error {
    if req.File == nil {
        return echo.NewHTTPError(http.StatusBadRequest, "file is required")
    }
    return nil
}