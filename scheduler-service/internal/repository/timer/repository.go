package timer_repository

import (
	"context"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/cowoking/scheduler-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type TimerRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *TimerRepository {
	return &TimerRepository{
		Postgres: pg,
	}
}

func (r *TimerRepository) Create(
	ctx context.Context,
	timer entity.Timer,
) (uuid.UUID, error) {

	query, args, _ := r.Builder.
		Insert("timer").
		Columns(
			"timer_type_id",
			"booking_id",
			"user_id",
			"trigger_at",
			"payload",
		).
		Values(
			timer.Type.ID,
			timer.BookingID,
			timer.UserID,
			timer.TriggerAt,
			timer.Payload,
		).
		Suffix("RETURNING id").
		ToSql()

	var id uuid.UUID

	err := r.GetTxManager(ctx).
		QueryRow(ctx, query, args...).
		Scan(&id)

	if err != nil {

		logrus.WithFields(logrus.Fields{
			"timer_type": timer.Type.Name,
			"booking_id": timer.BookingID.String(),
		}).WithError(err).Error("failed to create timer")

		return uuid.Nil, err
	}

	logrus.WithFields(logrus.Fields{
		"timer_id":   id.String(),
		"timer_type": timer.Type.Name,
		"booking_id": timer.BookingID.String(),
	}).Info("timer created")

	return id, nil
}

func (r *TimerRepository) CancelByBooking(
	ctx context.Context,
	bookingID uuid.UUID,
) error {

	query, args, _ := r.Builder.
		Update("timer").
		Set("status_id", 3). // cancelled
		Set("cancelled_at", time.Now()).
		Where("booking_id = ?", bookingID).
		Where("status_id = ?", 1).
		ToSql()

	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	if err != nil {

		logrus.WithField("booking_id", bookingID.String()).
			WithError(err).
			Error("failed to cancel timers")

		return err
	}

	logrus.WithFields(logrus.Fields{
		"booking_id": bookingID.String(),
		"affected":   cmd.RowsAffected(),
	}).Info("timers cancelled")

	return nil
}

func (r *TimerRepository) FindDueTimers(
	ctx context.Context,
	limit int,
) ([]entity.Timer, error) {

	query := `
		SELECT *
		FROM timer
		WHERE status_id = 1
		AND trigger_at <= NOW()
		ORDER BY trigger_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.GetTxManager(ctx).Query(ctx, query, limit)

	if err != nil {
		logrus.WithError(err).Error("failed to fetch due timers")
		return nil, err
	}

	defer rows.Close()

	rawTimers, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawTimer])

	if err != nil {
		logrus.WithError(err).Error("failed to collect timers")
		return nil, err
	}

	return lo.Map(rawTimers, func(r rawTimer, _ int) entity.Timer {
		return r.toEntity()
	}), nil
}

func (r *TimerRepository) MarkTriggered(
	ctx context.Context,
	timerIDs []uuid.UUID,
) error {

	query, args, _ := r.Builder.
		Update("timer").
		Set("status_id", 2).
		Set("triggered_at", time.Now()).
		Where("id = ANY(?)", timerIDs).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	if err != nil {

		logrus.WithFields(logrus.Fields{
			"count": len(timerIDs),
		}).WithError(err).
			Error("failed to mark timers triggered")

		return err
	}

	return nil
}
