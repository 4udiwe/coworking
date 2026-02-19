package booking_repository

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/cowoking/booking-service/internal/entity"
	. "github.com/4udiwe/cowoking/booking-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type BookingRepository struct {
	postgres.Postgres
}

func New(pg postgres.Postgres) *BookingRepository {
	return &BookingRepository{
		Postgres: pg,
	}
}

func (r *BookingRepository) Create(
	ctx context.Context,
	booking entity.Booking,
) error {

	query, args, _ := r.Builder.
		Insert("booking").
		Columns(
			"id",
			"user_id",
			"place_id",
			"start_time",
			"end_time",
			"status_id",
		).
		Values(
			booking.ID,
			booking.UserID,
			booking.Place.ID,
			booking.StartTime,
			booking.EndTime,
			StatusActive,
		).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		mapped := MapPgError(err)

		logrus.WithFields(logrus.Fields{
			"booking_id": booking.ID.String(),
			"place_id":   booking.Place.ID.String(),
			"user_id":    booking.UserID.String(),
		}).Warn("failed to create booking")
		return mapped
	}

	logrus.WithField("booking_id", booking.ID.String()).Info("booking created")

	return nil
}

func (r *BookingRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Booking, error) {

	query, args, _ := r.Builder.
		Select(
			"b.id",
			"b.user_id",
			"b.place_id",
			"p.label as place_label",
			"p.type as place_type",
			"p.coworking_id as place_coworking_id",
			"p.is_active as place_is_active",
			"b.start_time",
			"b.end_time",
			"b.status_id",
			"bs.name as status_name",
			"b.cancel_reason",
			"b.created_at",
			"b.updated_at",
			"b.cancelled_at",
		).
		From("booking b").
		Join("place p ON b.place_id = p.id").
		Join("booking_status bs ON b.status_id = bs.id").
		Where("b.id = ?", id).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).WithField("booking_id", id.String()).Error("failed to get booking")
		return entity.Booking{}, err
	}
	defer rows.Close()

	raw, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rawBookingPlaceStatus])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Booking{}, ErrBookingNotFound
		}
		logrus.WithError(err).WithField("booking_id", id.String()).Error("failed to get booking")
		return entity.Booking{}, err
	}

	return raw.toEntity(), nil
}

func (r *BookingRepository) ListByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Booking, error) {

	query, args, _ := r.Builder.
		Select(
			"b.id",
			"b.user_id",
			"b.place_id",
			"p.label as place_label",
			"p.type as place_type",
			"p.coworking_id as place_coworking_id",
			"p.is_active as place_is_active",
			"b.start_time",
			"b.end_time",
			"b.status_id",
			"bs.name as status_name",
			"b.cancel_reason",
			"b.created_at",
			"b.updated_at",
			"b.cancelled_at",
		).
		From("booking b").
		Join("place p ON b.place_id = p.id").
		Join("booking_status bs ON b.status_id = bs.id").
		Where("b.user_id = ?", userID).
		OrderBy("b.start_time DESC").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID.String()).Error("failed to list user bookings")
		return nil, err
	}
	defer rows.Close()

	raws, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawBookingPlaceStatus])

	if err != nil {
		logrus.WithError(err).Error("failed to collect booking rows")
		return nil, err
	}

	return lo.Map(raws, func(raw rawBookingPlaceStatus, _ int) entity.Booking {
		return raw.toEntity()
	}), nil
}

func (r *BookingRepository) Cancel(
	ctx context.Context,
	id uuid.UUID,
	reason *string,
) error {

	query, args, _ := r.Builder.
		Update("booking").
		Set("status_id", StatusCancelled).
		Set("cancel_reason", reason).
		Set("cancelled_at", time.Now()).
		Set("updated_at", time.Now()).
		Where("status_id = ?", StatusActive).
		Where("id = ?", id).
		ToSql()

	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("booking_id", id.String()).Error("failed to cancel booking")
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrBookingNotFound
	}

	logrus.WithField("booking_id", id.String()).Info("booking cancelled")

	return nil
}

func (r *BookingRepository) MarkCompleted(
	ctx context.Context,
	id uuid.UUID,
) error {

	query, args, _ := r.Builder.
		Update("booking").
		Set("status_id", StatusCompleted).
		Set("updated_at", time.Now()).
		Where("status_id = ?", StatusActive).
		Where("id = ?", id).
		ToSql()

	cmd, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("booking_id", id.String()).Error("failed to complete booking")
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrBookingNotFound
	}

	return nil
}
