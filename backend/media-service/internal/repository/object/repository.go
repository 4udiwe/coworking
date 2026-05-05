package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

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

func (s *Storage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)

	err := s.client.Upload(ctx, key, reader, int64(len(data)), contentType)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage upload failed")
		return err
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.Get(ctx, key)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage get failed")
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage read failed")
		return nil, err
	}

	return data, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	err := s.client.Delete(ctx, key)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage delete failed")
		return err
	}

	return nil
}

func (s *Storage) GeneratePresignedURL(
	ctx context.Context,
	key string,
	expiry time.Duration,
) (string, error) {

	url, err := s.client.GeneratePresignedURL(ctx, key, expiry)
	if err != nil {
		logrus.WithError(err).WithField("key", key).Error("storage presigned url failed")
		return "", err
	}

	return url, nil
}

func BuildKey(ownerID, mediaID, size string) string {
	return fmt.Sprintf("coworkings/%s/%s/%s.webp", ownerID, mediaID, size)
}
