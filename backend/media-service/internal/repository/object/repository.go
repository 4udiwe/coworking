package object_repository

import (
	"bytes"
	"context"
	"fmt"
	"io"

	minioClient "github.com/4udiwe/coworking/backend/media-service/pkg/minio"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	client *minioClient.Client
}

func New(client *minioClient.Client) *Storage {
	return &Storage{
		client: client,
	}
}

// Upload загружает файл из памяти ([]byte).
func (s *Storage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	// bytes.Reader реализует io.Reader, это lossless преобразование
	reader := bytes.NewReader(data)

	// Вызываем MinIO клиент с преобразованными параметрами
	err := s.client.Upload(ctx, key, reader, int64(len(data)), contentType)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage upload failed")
		return err
	}

	logrus.WithField("key", key).WithField("size", len(data)).Debug("file uploaded")
	return nil
}

// Get скачивает файл ([]byte).
// Адаптирует: *minio.Object → []byte (io.ReadAll) для медиа-сервиса.
func (s *Storage) Get(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.Get(ctx, key)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage get failed")
		return nil, err
	}
	defer obj.Close()

	// Читаем весь content
	data, err := io.ReadAll(obj)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage read failed")
		return nil, err
	}

	logrus.WithField("key", key).WithField("size", len(data)).Debug("file downloaded")
	return data, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	err := s.client.Delete(ctx, key)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage delete failed")
		return err
	}

	logrus.WithField("key", key).Debug("file deleted")
	return nil
}

func BuildKey(mediaID, size string) string {
    return fmt.Sprintf("/%s/%s.webp", mediaID, size)
}
