package booking_repository

import (
	"context"
	"errors"
	"fmt"
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
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *BookingRepository {
	return &BookingRepository{
		Postgres: pg,
	}
}

func (r *BookingRepository) Create(
	ctx context.Context,
	booking entity.Booking,
) (uuid.UUID, error) {

	query, args, _ := r.Builder.
		Insert("booking").
		Columns(
			"user_id",
			"user_name",
			"place_id",
			"start_time",
			"end_time",
			"status_id",
		).
		Values(
			booking.UserID,
			booking.UserName,
			booking.Place.ID,
			booking.StartTime,
			booking.EndTime,
			StatusActive,
		).
		Suffix("RETURNING id").
		ToSql()

	var id uuid.UUID

	err := r.GetTxManager(ctx).
		QueryRow(ctx, query, args...).
		Scan(&id)

	if err != nil {
		mapped := MapPgError(err)

		logrus.WithFields(logrus.Fields{
			"booking_id": id.String(),
			"place_id":   booking.Place.ID.String(),
			"user_id":    booking.UserID.String(),
		}).Warnf("failed to create booking: %v", err)

		return uuid.Nil, mapped
	}

	logrus.WithField("booking_id", id.String()).Info("booking created")

	return id, nil
}

func (r *BookingRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Booking, error) {

	query, args, _ := r.Builder.
		Select(
			"b.id",
			"b.user_id",
			"b.user_name",
			"b.place_id",
			"p.label as place_label",
			"p.place_type as place_type",
			"p.coworking_id as place_coworking_id",
			"p.is_active as place_is_active",
			"c.name as coworking_name",
			"c.address as coworking_address",
			"c.is_active as coworking_is_active",
			"c.created_at as coworking_created_at",
			"c.updated_at as coworking_updated_at",
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
		Join("coworking c ON p.coworking_id = c.id").
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
	page int,
	pageSize int,
	status *string,
) ([]entity.Booking, int, error) {

	// Get total count
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("booking b").
		Join("place p ON b.place_id = p.id").
		Join("booking_status bs ON b.status_id = bs.id").
		Where("b.user_id = ?", userID)

	if status != nil {
		countQuery = countQuery.Where("bs.name = ?", *status)
	}

	countSql, countArgs, _ := countQuery.ToSql()

	var totalCount int
	err := r.GetTxManager(ctx).QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID.String()).Error("failed to count user bookings")
		return nil, 0, err
	}

	// Get paginated bookings
	offset := (page - 1) * pageSize
	query := r.Builder.
		Select(
			"b.id",
			"b.user_id",
			"b.user_name",
			"b.place_id",
			"p.label as place_label",
			"p.place_type as place_type",
			"p.coworking_id as place_coworking_id",
			"p.is_active as place_is_active",
			"c.name as coworking_name",
			"c.address as coworking_address",
			"c.is_active as coworking_is_active",
			"c.created_at as coworking_created_at",
			"c.updated_at as coworking_updated_at",
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
		Join("coworking c ON p.coworking_id = c.id").
		Join("booking_status bs ON b.status_id = bs.id").
		Where("b.user_id = ?", userID)

	if status != nil {
		query = query.Where("bs.name = ?", *status)
	}

	query = query.
		OrderBy(fmt.Sprintf("CASE WHEN b.status_id = %d THEN 0 ELSE 1 END", StatusActive)).
		OrderBy("b.start_time DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset))

	sql, args, _ := query.ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, sql, args...)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID.String()).Error("failed to list user bookings")
		return nil, 0, err
	}
	defer rows.Close()

	raws, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawBookingPlaceStatus])

	if err != nil {
		logrus.WithError(err).Error("failed to collect booking rows")
		return nil, 0, err
	}

	bookings := lo.Map(raws, func(raw rawBookingPlaceStatus, _ int) entity.Booking {
		return raw.toEntity()
	})

	return bookings, totalCount, nil
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

func (r *BookingRepository) GetAdminActiveBookings(
	ctx context.Context,
	coworkingID uuid.UUID,
	page int,
	pageSize int,
	dateFrom *time.Time,
	dateTo *time.Time,
	placeType *string,
	sortBy *string,
) ([]entity.Booking, int, error) {

	// Get total count
	countQuery := r.Builder.
		Select("COUNT(*)").
		From("booking b").
		Join("place p ON b.place_id = p.id").
		Join("booking_status bs ON b.status_id = bs.id").
		Where("p.coworking_id = ?", coworkingID).
		Where("bs.name = ?", entity.BookingStatusActive)

	// Apply filters to count query
	if dateFrom != nil {
		countQuery = countQuery.Where("b.start_time >= ?", *dateFrom)
	}
	if dateTo != nil {
		countQuery = countQuery.Where("b.start_time <= ?", *dateTo)
	}
	if placeType != nil {
		countQuery = countQuery.Where("p.place_type = ?", *placeType)
	}

	countSql, countArgs, _ := countQuery.ToSql()

	var totalCount int
	err := r.GetTxManager(ctx).QueryRow(ctx, countSql, countArgs...).Scan(&totalCount)
	if err != nil {
		logrus.WithError(err).WithField("coworking_id", coworkingID.String()).Error("failed to count admin active bookings")
		return nil, 0, err
	}

	// Get paginated bookings
	offset := (page - 1) * pageSize
	query := r.Builder.
		Select(
			"b.id",
			"b.user_id",
			"b.user_name",
			"b.place_id",
			"p.label as place_label",
			"p.place_type as place_type",
			"p.coworking_id as place_coworking_id",
			"p.is_active as place_is_active",
			"c.name as coworking_name",
			"c.address as coworking_address",
			"c.is_active as coworking_is_active",
			"c.created_at as coworking_created_at",
			"c.updated_at as coworking_updated_at",
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
		Join("coworking c ON p.coworking_id = c.id").
		Join("booking_status bs ON b.status_id = bs.id").
		Where("p.coworking_id = ?", coworkingID).
		Where("bs.name = ?", entity.BookingStatusActive)

	// Apply filters
	if dateFrom != nil {
		query = query.Where("b.start_time >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("b.start_time <= ?", *dateTo)
	}
	if placeType != nil {
		query = query.Where("p.place_type = ?", *placeType)
	}

	// Apply sorting
	orderBy := "b.start_time DESC"
	if sortBy != nil {
		if *sortBy == "asc" {
			orderBy = "b.start_time ASC"
		}
	}

	query = query.
		OrderBy(orderBy).
		Limit(uint64(pageSize)).
		Offset(uint64(offset))

	sql, args, _ := query.ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, sql, args...)
	if err != nil {
		logrus.WithError(err).WithField("coworking_id", coworkingID.String()).Error("failed to list admin active bookings")
		return nil, 0, err
	}
	defer rows.Close()

	raws, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawBookingPlaceStatus])

	if err != nil {
		logrus.WithError(err).Error("failed to collect booking rows")
		return nil, 0, err
	}

	bookings := lo.Map(raws, func(raw rawBookingPlaceStatus, _ int) entity.Booking {
		return raw.toEntity()
	})

	return bookings, totalCount, nil
}
