package place_repository

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	. "github.com/4udiwe/cowoking/booking-service/internal/repository"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type PlaceRepository struct {
	postgres.Postgres
}

func New(pg postgres.Postgres) *PlaceRepository {
	return &PlaceRepository{
		Postgres: pg,
	}
}

func (r *PlaceRepository) CreateBatch(
	ctx context.Context,
	places []entity.Place,
) error {

	if len(places) == 0 {
		return nil
	}

	builder := r.Builder.
		Insert("place").
		Columns(
			"coworking_id",
			"label",
			"place_type",
		)

	for _, place := range places {
		builder = builder.Values(
			place.Coworking.ID,
			place.Label,
			place.PlaceType,
		)
	}

	query, args, _ := builder.ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)

		logrus.Error("failed to create places batch: ", mapped)
		return mapped
	}

	return nil
}

func (r *PlaceRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Place, error) {

	query, args, _ := r.Builder.
		Select(
			"p.id",
			"p.label",
			"p.place_type",
			"p.is_active",
			"c.id AS coworking_id",
			"c.name AS coworking_name",
			"c.address AS coworking_address",
			"c.is_active AS coworking_is_active",
		).
		From("place p").
		Join("coworking ON coworking.id = p.coworking_id c").
		Where("p.id = ?", id).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).WithField("place_id", id.String()).Error("failed to get place")
		return entity.Place{}, err
	}
	defer rows.Close()

	raw, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rawPlaceCoworking])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Place{}, ErrPlaceNotFound
		}
		logrus.WithError(err).WithField("place_id", id.String()).Error("failed to get place")
		return entity.Place{}, err
	}

	return raw.toEntity(), nil
}

func (r *PlaceRepository) GetByCoworking(
	ctx context.Context,
	coworkingID uuid.UUID,
) ([]entity.Place, error) {

	query, args, _ := r.Builder.
		Select(
			"p.id",
			"p.label",
			"p.place_type",
			"p.is_active",
			"c.id AS coworking_id",
			"c.name AS coworking_name",
			"c.address AS coworking_address",
			"c.is_active AS coworking_is_active",
		).
		From("place p").
		Join("coworking ON coworking.id = p.coworking_id c").
		Where("p.coworking_id = ?", coworkingID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).WithField("coworking_id", coworkingID.String()).Error("failed to get places by coworking")
		return nil, err
	}
	defer rows.Close()

	rawPlaces, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawPlaceCoworking])

	if err != nil {
		logrus.WithError(err).WithField("coworking_id", coworkingID.String()).Error("failed to get places by coworking")
		return nil, err
	}

	return lo.Map(rawPlaces, func(r rawPlaceCoworking, _ int) entity.Place {
		return r.toEntity()
	}), nil
}

func (r *PlaceRepository) GetAvailableByCoworking(
	ctx context.Context,
	coworkingID uuid.UUID,
	start time.Time,
	end time.Time,
) ([]entity.Place, error) {

	query, args, _ := r.Builder.
		Select(
			"p.id",
			"p.label",
			"p.place_type",
			"p.is_active",
			"c.id AS coworking_id",
			"c.name AS coworking_name",
			"c.address AS coworking_address",
			"c.is_active AS coworking_is_active",
		).
		From("place p").
		Join("coworking ON coworking.id = p.coworking_id c").
		Where("p.coworking_id = ?", coworkingID).
		Where("p.is_active = TRUE").
		Where(`p.id NOT IN (
			SELECT place_id FROM booking
			WHERE status_id = 1 -- active
			AND start_time < ?
			AND end_time > ?
		)`, end, start).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).WithField("coworking_id", coworkingID.String()).Error("failed to get available places by coworking")
		return nil, err
	}
	defer rows.Close()

	rawPlaces, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawPlaceCoworking])

	if err != nil {
		logrus.WithError(err).WithField("coworking_id", coworkingID.String()).Error("failed to get available places by coworking")
		return nil, err
	}

	return lo.Map(rawPlaces, func(r rawPlaceCoworking, _ int) entity.Place {
		return r.toEntity()
	}), nil
}

func (r *PlaceRepository) SetActive(
	ctx context.Context,
	id uuid.UUID,
	active bool,
) error {

	query, args, _ := r.Builder.
		Update("place").
		Set("is_active", active).
		Where("id = ?", id).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("place_id", id.String()).Error("failed to set place active")
		return err
	}
	return nil
}

func (r *PlaceRepository) CheckHasActiveBookings(
	ctx context.Context,
	placeID uuid.UUID,
) (bool, error) {

	query, args, _ := r.Builder.
		Select("1").
		From("booking").
		Where("place_id = ?", placeID).
		Where("status_id = ?", squirrel.Expr("SELECT id FROM booking_status WHERE name = ?", entity.BookingStatusActive)).
		Limit(1).
		ToSql()

	var hasActive bool

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&hasActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		logrus.WithError(err).WithField("place_id", placeID.String()).Error("failed to check active bookings for place")
		return false, err
	}

	return hasActive, nil
}