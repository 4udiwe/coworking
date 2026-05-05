package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/coworking/backend/media-service/internal/entity"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "media"

var (
	ErrMediaNotFound = errors.New("media not found")
)

// MediaRepository — реализация поверх MongoDB.
type MediaRepository struct {
	coll *mongo.Collection
}

func New(db *mongo.Database) *MediaRepository {
	return &MediaRepository{
		coll: db.Collection(collectionName),
	}
}

func (r *MediaRepository) Create(ctx context.Context, media entity.Media) (primitive.ObjectID, error) {
	media.CreatedAt = time.Now()
	media.UpdatedAt = time.Now()

	result, err := r.coll.InsertOne(ctx, media)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"owner_type": media.OwnerType,
			"owner_id":   media.OwnerID,
			"purpose":    media.Purpose,
		}).Error("failed to create media")
		return primitive.NilObjectID, err
	}

	id := result.InsertedID.(primitive.ObjectID)
	logrus.WithField("media_id", id.Hex()).Info("media created")

	return id, nil
}

func (r *MediaRepository) AddVariant(ctx context.Context, id primitive.ObjectID, variant entity.ImageVariant) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$push": bson.M{"variants": variant},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithField("media_id", id.Hex()).Error("failed to add variant")
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMediaNotFound
	}

	return nil
}

func (r *MediaRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status entity.ProcessingStatus) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"media_id": id.Hex(),
			"status":   status,
		}).Error("failed to update media status")
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMediaNotFound
	}

	return nil
}

func (r *MediaRepository) UpdateStatusAndVariants(
	ctx context.Context,
	id primitive.ObjectID,
	status entity.ProcessingStatus,
	variants []entity.ImageVariant,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
		"$push": bson.M{
			"variants": bson.M{"$each": variants},
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithField("media_id", id.Hex()).Error("failed to update status and variants")
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMediaNotFound
	}

	logrus.WithFields(logrus.Fields{
		"media_id":       id.Hex(),
		"status":         status,
		"variants_added": len(variants),
	}).Info("media processing complete")

	return nil
}

func (r *MediaRepository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	filter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false}, // идемпотентность: уже удалённые не трогаем
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithField("media_id", id.Hex()).Error("failed to soft delete media")
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMediaNotFound
	}

	logrus.WithField("media_id", id.Hex()).Info("media soft deleted")

	return nil
}

func (r *MediaRepository) SoftDeleteCoverByOwner(ctx context.Context, ownerType, ownerID string) error {
	now := time.Now()
	filter := bson.M{
		"owner_type": ownerType,
		"owner_id":   ownerID,
		"purpose":    entity.PurposeCover,
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	// UpdateMany на случай если по какой-то причине обложек несколько (defensive)
	_, err := r.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"owner_type": ownerType,
			"owner_id":   ownerID,
		}).Error("failed to soft delete cover")
		return err
	}

	return nil
}

func (r *MediaRepository) UpdateSortOrder(ctx context.Context, orders map[primitive.ObjectID]int) error {
	// Используем bulk write для атомарного обновления всех порядков за один round-trip
	models := make([]mongo.WriteModel, 0, len(orders))

	for id, order := range orders {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"sort_order": order,
					"updated_at": time.Now(),
				},
			}),
		)
	}

	_, err := r.coll.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		logrus.WithError(err).Error("failed to update sort orders")
		return err
	}

	return nil
}

func (r *MediaRepository) IncrementRetryCount(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"retry_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithField("media_id", id.Hex()).Error("failed to increment retry count")
		return err
	}

	return nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id primitive.ObjectID) (entity.Media, error) {
	filter := bson.M{"_id": id}

	var media entity.Media
	err := r.coll.FindOne(ctx, filter).Decode(&media)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Media{}, ErrMediaNotFound
		}
		logrus.WithError(err).WithField("media_id", id.Hex()).Error("failed to get media by id")
		return entity.Media{}, err
	}

	return media, nil
}

func (r *MediaRepository) GetByOwner(ctx context.Context, ownerType, ownerID string) ([]entity.Media, error) {
	filter := bson.M{
		"owner_type": ownerType,
		"owner_id":   ownerID,
		"deleted_at": bson.M{"$exists": false},
	}
	opts := options.Find().SetSort(bson.D{
		{Key: "sort_order", Value: 1},
		{Key: "created_at", Value: 1},
	})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"owner_type": ownerType,
			"owner_id":   ownerID,
		}).Error("failed to get media by owner")
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []entity.Media
	if err := cursor.All(ctx, &result); err != nil {
		logrus.WithError(err).Error("failed to decode media list")
		return nil, err
	}

	return result, nil
}

func (r *MediaRepository) GetCoverByOwner(ctx context.Context, ownerType, ownerID string) (entity.Media, error) {
	filter := bson.M{
		"owner_type": ownerType,
		"owner_id":   ownerID,
		"purpose":    entity.PurposeCover,
		"deleted_at": bson.M{"$exists": false},
	}

	var media entity.Media
	err := r.coll.FindOne(ctx, filter).Decode(&media)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Media{}, ErrMediaNotFound
		}
		logrus.WithError(err).WithFields(logrus.Fields{
			"owner_type": ownerType,
			"owner_id":   ownerID,
		}).Error("failed to get cover by owner")
		return entity.Media{}, err
	}

	return media, nil
}

func (r *MediaRepository) CountGalleryByOwner(ctx context.Context, ownerType, ownerID string) (int, error) {
	filter := bson.M{
		"owner_type": ownerType,
		"owner_id":   ownerID,
		"purpose":    entity.PurposeGallery,
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"owner_type": ownerType,
			"owner_id":   ownerID,
		}).Error("failed to count gallery media")
		return 0, err
	}

	return int(count), nil
}

func (r *MediaRepository) FindStale(ctx context.Context, threshold time.Duration, limit int) ([]entity.Media, error) {
	cutoff := time.Now().Add(-threshold)

	filter := bson.M{
		"status":     entity.StatusProcessing,
		"updated_at": bson.M{"$lt": cutoff},
	}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "updated_at", Value: 1}}) // самые старые первыми

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).Error("failed to find stale media")
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []entity.Media
	if err := cursor.All(ctx, &result); err != nil {
		logrus.WithError(err).Error("failed to decode stale media list")
		return nil, err
	}

	return result, nil
}

func (r *MediaRepository) GetByIDs(
	ctx context.Context,
	ids []primitive.ObjectID,
) ([]entity.Media, error) {

	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("failed to get media by ids")
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []entity.Media
	if err := cursor.All(ctx, &result); err != nil {
		logrus.WithError(err).Error("failed to decode media list")
		return nil, err
	}

	return result, nil
}

func (r *MediaRepository) UpdatePurpose(ctx context.Context, id primitive.ObjectID, purpose entity.MediaPurpose) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"purpose":    purpose,
			"updated_at": time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"media_id": id.Hex(),
			"purpose":  purpose,
		}).Error("failed to update media purpose")
		return err
	}

	if result.MatchedCount == 0 {
		return ErrMediaNotFound
	}

	return nil
}