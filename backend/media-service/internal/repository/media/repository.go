package media_repository

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
	media.Variants = []entity.ImageVariant{}

	result, err := r.coll.InsertOne(ctx, media)
	if err != nil {
		logrus.WithError(err).Error("failed to create media")
		return primitive.NilObjectID, err
	}

	id := result.InsertedID.(primitive.ObjectID)
	logrus.WithField("media_id", id.Hex()).Info("media created")

	return id, nil
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

func (r *MediaRepository) Delete(
	ctx context.Context,
	id primitive.ObjectID,
) error {
	filter := bson.M{"_id": id}
	result, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		logrus.WithError(err).WithField("media_id", id.Hex()).Error("failed to delete media")
		return err
	}
	if result.DeletedCount == 0 {
		return ErrMediaNotFound
	}
	logrus.WithField("media_id", id.Hex()).Info("media deleted")
	return nil
}

func (r *MediaRepository) FindExpired(
	ctx context.Context,
	now time.Time,
	limit int,
) ([]entity.Media, error) {
	filter := bson.M{
		//"status": entity.StatusProcessed,
		"variants": bson.M{
			"$elemMatch": bson.M{
				"expires_at": bson.M{"$lt": now},
			},
		},
	}
	opts := options.Find().SetLimit(int64(limit))
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).Error("failed to find expired media")
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []entity.Media
	if err := cursor.All(ctx, &result); err != nil {
		logrus.WithError(err).Error("failed to decode expired media list")
		return nil, err
	}
	return result, nil
}
