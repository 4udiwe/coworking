package device_repository

import (
	"context"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type DeviceRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *DeviceRepository {
	return &DeviceRepository{
		Postgres: pg,
	}
}

func (r *DeviceRepository) Create(
	ctx context.Context,
	device entity.UserDevice,
) (uuid.UUID, error) {

	query, args, _ := r.Builder.
		Insert("user_device").
		Columns(
			"user_id",
			"device_token",
			"platform",
		).
		Values(
			device.UserID,
			device.DeviceToken,
			device.Platform,
		).
		Suffix("RETURNING id").
		ToSql()

	var id uuid.UUID

	err := r.GetTxManager(ctx).
		QueryRow(ctx, query, args...).
		Scan(&id)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id":  device.UserID.String(),
			"platform": device.Platform,
		}).WithError(err).Error("failed to create device")
		return uuid.Nil, err
	}

	return id, nil
}

func (r *DeviceRepository) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.UserDevice, error) {

	query, args, _ := r.Builder.
		Select("*").
		From("user_device").
		Where("user_id = ?", userID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).Error("failed to fetch devices")
		return nil, err
	}

	defer rows.Close()

	rawDevices, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawDevice])

	if err != nil {
		logrus.WithError(err).Error("failed to collect devices")
		return nil, err
	}

	return lo.Map(rawDevices, func(r rawDevice, _ int) entity.UserDevice {
		return r.toEntity()
	}), nil
}

func (r *DeviceRepository) DeleteByToken(ctx context.Context, token string) error {
	query, args, _ := r.Builder.
		Delete("user_device").
		Where("device_token = ?", token).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}
