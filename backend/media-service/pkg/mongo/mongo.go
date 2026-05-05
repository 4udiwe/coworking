package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultConnAttempts   = 5
	defaultConnRetryDelay = 2 * time.Second
	defaultPingTimeout    = 5 * time.Second

	collectionName = "media"
)

// MongoDB — обёртка над mongo.Client.
type MongoDB struct {
	connTimeout  time.Duration
	connAttempts int

	Client   *mongo.Client
	Database *mongo.Database
}

// Cоздаёт подключение к MongoDB с retry-логикой.
func New(uri, dbName string, opts ...Option) (*MongoDB, error) {
	m := &MongoDB{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnectTimeout,
	}

	for _, opt := range opts {
		opt(m)
	}

	clientOpts := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(m.connTimeout).
		SetServerSelectionTimeout(m.connTimeout)

	var (
		client *mongo.Client
		err    error
	)

	// Retry loop
	for attempt := 1; attempt <= m.connAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), m.connTimeout)

		client, err = mongo.Connect(ctx, clientOpts)
		cancel()

		if err != nil {
			logrus.Warnf("mongodb - New - attempt %d/%d connect: %v", attempt, m.connAttempts, err)
			time.Sleep(defaultConnRetryDelay)
			continue
		}

		// Ping для проверки живого соединения
		pingCtx, pingCancel := context.WithTimeout(context.Background(), defaultPingTimeout)
		err = client.Ping(pingCtx, readpref.Primary())
		pingCancel()

		if err != nil {
			logrus.Warnf("mongodb - New - attempt %d/%d ping: %v", attempt, m.connAttempts, err)
			_ = client.Disconnect(context.Background())
			time.Sleep(defaultConnRetryDelay)
			continue
		}

		logrus.Infof("mongodb - New - connected on attempt %d", attempt)
		break
	}

	if err != nil {
		return nil, fmt.Errorf("mongodb - New - failed after %d attempts: %w", m.connAttempts, err)
	}

	m.Client = client
	m.Database = client.Database(dbName)

	// Создаём индексы сразу при старте. Идемпотентно — безопасно вызывать повторно.
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer idxCancel()

	if err := ensureIndexes(idxCtx, m.Database); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongodb - New - ensureIndexes: %w", err)
	}

	return m, nil
}

// Close корректно закрывает соединение.
// Вызывать через defer в app.Start().
func (m *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		logrus.WithError(err).Error("mongodb - Close - disconnect error")
	}
}

// Ping проверяет живость соединения.
// Используется в /readyz health check.
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, readpref.Primary())
}

// ensureIndexes создаёт индексы для коллекции media.
// Идемпотентно — если индекс уже существует, MongoDB просто пропускает его.
// Вызывается автоматически при инициализации клиента.
func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(collectionName)

	indexes := []mongo.IndexModel{
		// Основной запрос: "все активные фото коворкинга X", отсортированные для галереи.
		// Покрывает GetByOwner().
		{
			Keys: bson.D{
				{Key: "owner_type", Value: 1},
				{Key: "owner_id", Value: 1},
				{Key: "deleted_at", Value: 1},
				{Key: "sort_order", Value: 1},
			},
			Options: options.Index().SetName("idx_owner_active_sorted"),
		},
		// Запрос обложки коворкинга — самый частый, в каждой карточке.
		// Покрывает GetCoverByOwner() и SoftDeleteCoverByOwner().
		{
			Keys: bson.D{
				{Key: "owner_type", Value: 1},
				{Key: "owner_id", Value: 1},
				{Key: "purpose", Value: 1},
				{Key: "deleted_at", Value: 1},
			},
			Options: options.Index().SetName("idx_owner_purpose"),
		},
		// Stale checker: поиск застрявших в processing.
		// Partial index — маленький, работает только по status=processing.
		// Покрывает FindStale().
		{
			Keys: bson.D{
				{Key: "updated_at", Value: 1},
			},
			Options: options.Index().
				SetName("idx_stale_processing").
				SetPartialFilterExpression(bson.M{
					"status": string(entity.StatusProcessing),
				}),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("ensureIndexes: %w", err)
	}

	logrus.Info("mongodb - ensureIndexes - indexes created successfully")
	return nil
}
