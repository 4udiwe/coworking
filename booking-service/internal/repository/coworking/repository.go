package coworking_repository

import (
	"context"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	. "github.com/4udiwe/cowoking/booking-service/internal/repository"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type CoworkingRepository struct {
	postgres.Postgres
}

func New(pg postgres.Postgres) *CoworkingRepository {
	return &CoworkingRepository{
		Postgres: pg,
	}
}

func (r *CoworkingRepository) Create(ctx context.Context, coworking entity.Coworking) error {
	query, args, _ := r.Builder.
		Insert("coworking").
		Columns(
			"name",
			"address",
			"is_active",
		).
		Values(
			coworking.Name,
			coworking.Address,
			coworking.IsActive,
		).
		Suffix("RETURNING id").
		ToSql()

	var id uuid.UUID
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		mapped := MapPgError(err)
		logrus.Error("failed to create coworking: ", mapped)
		return mapped
	}

	return nil
}

func (r *CoworkingRepository) GetByID(ctx context.Context, id uuid.UUID) (entity.Coworking, error) {
	query, args, _ := r.Builder.
		Select(
			"id",
			"name",
			"address",
			"is_active",
			"created_at",
			"updated_at",
		).
		From("coworking").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	row, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", id.String()).Error("failed to get coworking by id: ", mapped)
		return entity.Coworking{}, mapped
	}
	defer row.Close()

	rawCoworking, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[rawCoworking])
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", id.String()).Error("failed to collect coworking row: ", mapped)
		return entity.Coworking{}, mapped
	}

	return rawCoworking.toEntity(), nil
}

func (r *CoworkingRepository) Update(ctx context.Context, coworking entity.Coworking) error {
	query, args, _ := r.Builder.
		Update("coworking").
		Set("name", coworking.Name).
		Set("address", coworking.Address).
		Set("is_active", coworking.IsActive).
		Where(squirrel.Eq{"id": coworking.ID}).
		ToSql()

	cmdTag, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworking.ID.String()).Error("failed to update coworking: ", mapped)
		return mapped
	}
	if cmdTag.RowsAffected() == 0 {
		logrus.WithField("coworking_id", coworking.ID.String()).Warn("no coworking updated")
		return ErrCoworkingNotFound
	}

	return nil
}

func (r *CoworkingRepository) List(ctx context.Context) ([]entity.Coworking, error) {
	query, args, _ := r.Builder.
		Select(
			"id",
			"name",
			"address",
			"is_active",
			"created_at",
			"updated_at",
		).
		From("coworking").
		OrderBy("created_at DESC").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.Error("failed to list coworkings: ", mapped)
		return nil, mapped
	}

	rawCoworkings, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawCoworking])
	if err != nil {
		mapped := MapPgError(err)
		logrus.Error("failed to collect coworkings: ", mapped)
		return nil, mapped
	}

	return lo.Map(rawCoworkings, func(c rawCoworking, _ int) entity.Coworking {
		return c.toEntity()
	}), nil
}

func (r *CoworkingRepository) CreateLayoutVersion(ctx context.Context, layout entity.CoworkingLayout) error {
	query, args, _ := r.Builder.
		Insert("coworking_layout").
		Columns(
			"coworking_id",
			"layout",
			"version",
		).
		Values(
			layout.CoworkingID,
			layout.Layout,
			layout.Version,
		).
		Suffix("RETURNING id").
		ToSql()

	cmdTag, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", layout.CoworkingID.String()).Error("failed to create coworking layout: ", mapped)
		return mapped
	}
	if cmdTag.RowsAffected() == 0 {
		logrus.WithField("coworking_id", layout.CoworkingID.String()).Warn("no coworking layout created")
		return ErrCoworkingNotFound
	}

	return nil
}

func (r *CoworkingRepository) GetLatestLayout(ctx context.Context, coworkingID uuid.UUID) (entity.CoworkingLayout, error) {
	query, args, _ := r.Builder.
		Select(
			"id",
			"coworking_id",
			"layout",
			"version",
			"created_at",
		).
		From("coworking_layout").
		Where(squirrel.Eq{"coworking_id": coworkingID}).
		OrderBy("version DESC").
		Limit(1).
		ToSql()

	row, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworkingID.String()).Error("failed to get latest layout: ", mapped)
		return entity.CoworkingLayout{}, mapped
	}
	defer row.Close()

	rawLayout, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[rawCoworkingLayout])
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworkingID.String()).Error("failed to collect layout row: ", mapped)
		return entity.CoworkingLayout{}, mapped
	}

	return rawLayout.toEntity(), nil
}

func (r *CoworkingRepository) GetLayoutByVersion(
	ctx context.Context,
	coworkingID uuid.UUID,
	version int,
) (entity.CoworkingLayout, error) {
	query, args, _ := r.Builder.
		Select(
			"id",
			"coworking_id",
			"layout",
			"version",
			"created_at",
		).
		From("coworking_layout").
		Where(squirrel.Eq{"coworking_id": coworkingID, "version": version}).
		ToSql()

	row, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithFields(logrus.Fields{
			"coworking_id": coworkingID.String(),
			"version":      version,
		}).Error("failed to get layout by version: ", mapped)
		return entity.CoworkingLayout{}, mapped
	}
	defer row.Close()

	rawLayout, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[rawCoworkingLayout])
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithFields(logrus.Fields{
			"coworking_id": coworkingID.String(),
			"version":      version,
		}).Error("failed to collect layout row: ", mapped)
		return entity.CoworkingLayout{}, mapped
	}

	return rawLayout.toEntity(), nil
}

func (r *CoworkingRepository) ListLayoutVersions(
	ctx context.Context,
	coworkingID uuid.UUID,
) ([]entity.CoworkingLayoutVersionTime, error) {
	query, args, _ := r.Builder.
		Select("version", "created_at").
		From("coworking_layout").
		Where(squirrel.Eq{"coworking_id": coworkingID}).
		OrderBy("version DESC").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworkingID.String()).Error("failed to list layout versions: ", mapped)
		return nil, mapped
	}
	defer rows.Close()

	rawLayoutVersions, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawLayoutVersionTime])
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworkingID.String()).Error("failed to collect layout versions: ", mapped)
		return nil, mapped
	}

	return lo.Map(rawLayoutVersions, func(v rawLayoutVersionTime, _ int) entity.CoworkingLayoutVersionTime {
		return v.toEntity()
	}), nil
}

func (r *CoworkingRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	logrus.WithFields(logrus.Fields{
		"coworking_id": id.String(),
		"active":       active,
	}).Info("setting coworking active status")

	query, args, _ := r.Builder.
		Update("coworking").
		Set("is_active", active).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	cmdTag, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", id.String()).Error("failed to set coworking active status: ", mapped)
		return mapped
	}
	if cmdTag.RowsAffected() == 0 {
		logrus.WithField("coworking_id", id.String()).Warn("no coworking updated")
		return ErrCoworkingNotFound
	}

	return nil
}

func (r *CoworkingRepository) RollbackLatestLayoutVersion(ctx context.Context, coworkingID uuid.UUID) error {
	logrus.WithField("coworking_id", coworkingID.String()).Info("rolling back latest layout version")

	query, args, _ := r.Builder.
		Delete("coworking_layout").
		Where("id = (SELECT id FROM coworking_layout WHERE coworking_id = ? ORDER BY version DESC LIMIT 1)", coworkingID).
		ToSql()

	cmdTag, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworkingID.String()).Error("failed to rollback latest layout version: ", mapped)
		return mapped
	}

	if cmdTag.RowsAffected() == 0 {
		logrus.WithField("coworking_id", coworkingID.String()).Warn("no layout version rolled back")
		return ErrLayoutNotFound
	}

	return nil
}

func (r *CoworkingRepository) CheckHasActiveBookings(ctx context.Context, coworkingID uuid.UUID) (bool, error) {
	query, args, _ := r.Builder.
		Select("1").
		From("booking b").
		Join("place p ON b.place_id = p.id").
		Where(squirrel.Eq{"p.coworking_id": coworkingID, "b.status": "active"}).
		Limit(1).
		ToSql()

	row, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)
		logrus.WithField("coworking_id", coworkingID.String()).Error("failed to check active bookings: ", mapped)
		return false, mapped
	}

	defer row.Close()

	hasActive := row.Next()

	return hasActive, nil
}
