package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultConnAttempts   = 5
	defaultRetryDelay     = 2 * time.Second
)

type Client struct {
	client *minio.Client
	bucket string

	publicEndpoint string

	connAttempts int
	connTimeout  time.Duration
}

func New(endpoint, publicEndpoint, accessKey, secretKey, bucket string, useSSL bool) (*Client, error) {
	c := &Client{
		connAttempts:   defaultConnAttempts,
		connTimeout:    defaultConnectTimeout,
		bucket:         bucket,
		publicEndpoint: publicEndpoint,
	}

	var err error

	for attempt := 1; attempt <= c.connAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), c.connTimeout)
		defer cancel()

		c.client, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})

		if err != nil {
			logrus.Warnf("minio - New - attempt %d/%d: %v", attempt, c.connAttempts, err)
			time.Sleep(defaultRetryDelay)
			continue
		}

		// Проверка доступа к bucket
		exists, err := c.client.BucketExists(ctx, bucket)
		if err != nil {
			logrus.Warnf("minio - New - bucket check failed: %v", err)
			time.Sleep(defaultRetryDelay)
			continue
		}

		if !exists {
			err = c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				logrus.Warnf("minio - New - make bucket failed: %v", err)
				time.Sleep(defaultRetryDelay)
				continue
			}
			logrus.Infof("minio - bucket %s created", bucket)
		}

		logrus.Infof("minio - connected on attempt %d", attempt)
		return c, nil
	}

	return nil, fmt.Errorf("minio - New - failed after %d attempts: %w", c.connAttempts, err)
}

func (c *Client) Upload(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) error {

	_, err := c.client.PutObject(ctx, c.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		logrus.WithError(err).WithField("object", objectName).Error("minio upload failed")
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, objectName string) (*minio.Object, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		logrus.WithError(err).WithField("object", objectName).Error("minio get failed")
		return nil, err
	}
	return obj, nil
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		logrus.WithError(err).WithField("object", objectName).Error("minio delete failed")
		return err
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.ListBuckets(ctx)
	return err
}

func (c *Client) GeneratePresignedURL(
	ctx context.Context,
	objectName string,
	expiry time.Duration,
) (string, error) {

	reqParams := make(url.Values)

	presignedURL, err := c.client.PresignedGetObject(
		ctx,
		c.bucket,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil {
		logrus.WithError(err).
			WithField("object", objectName).
			Error("minio presigned url failed")
		return "", err
	}

	u, err := url.Parse(presignedURL.String())
	if err != nil {
		return "", err
	}

	u.Host = c.publicEndpoint

	return u.String(), nil
}
