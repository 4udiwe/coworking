package notification_repository

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/coworking/notification-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *NotificationRepository {
	return &NotificationRepository{
		Postgres: pg,
	}
}

func (r *NotificationRepository) Create(
	ctx context.Context,
	notification entity.Notification,
) (uuid.UUID, error) {

	query := `
		INSERT INTO notification (
			user_id,
			notification_type_id,
			title,
			body,
			payload,
			status_id
		) VALUES (
			$1,
			(SELECT id FROM notification_type WHERE name = $2),
			$3,
			$4,
			$5,
			1 -- unread
		)
		RETURNING id
	`

	var id uuid.UUID
	err := r.GetTxManager(ctx).QueryRow(
		ctx,
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Body,
		notification.Payload,
	).Scan(&id)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": notification.UserID.String(),
			"type":    notification.Type,
		}).WithError(err).Error("failed to create notification")
		return uuid.Nil, err
	}

	return id, nil
}

func (r *NotificationRepository) FindByUser(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]entity.Notification, error) {

	query, args, _ := r.Builder.
		Select(
			"n.id",
			"n.user_id",
			"n.notification_type_id",
			"nt.name as notification_type_name",
			"n.title",
			"n.body",
			"n.payload",
			"n.status_id",
			"ns.name as status_name",
			"n.created_at",
			"n.read_at",
		).
		From("notification n").
		Join("notification_type nt ON n.notification_type_id = nt.id").
		Join("notification_status ns ON n.status_id = ns.id").
		Where("n.user_id = ?", userID).
		OrderBy("n.created_at DESC").
		Limit(uint64(limit)).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).Error("failed to fetch notifications")
		return nil, err
	}

	defer rows.Close()

	rawNotifications, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawNotification])

	if err != nil {
		logrus.WithError(err).Error("failed to collect notifications")
		return nil, err
	}

	return lo.Map(rawNotifications, func(r rawNotification, _ int) entity.Notification {
		return r.toEntity()
	}), nil
}
func (r *NotificationRepository) FetchUnreadByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.Notification, error) {
	query, args, _ := r.Builder.
		Select(
			"n.id",
			"n.user_id",
			"n.notification_type_id",
			"nt.name as notification_type_name",
			"n.title",
			"n.body",
			"n.payload",
			"n.status_id",
			"ns.name as status_name",
			"n.created_at",
			"n.read_at",
		).
		From("notification n").
		Join("notification_type nt ON n.notification_type_id = nt.id").
		Join("notification_status ns ON n.status_id = ns.id").
		Where("n.user_id = ?", userID).
		Where("ns.status_name = ?", entity.StatusUnread).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).Error("failed to fetch unread notifications")
		return nil, err
	}

	defer rows.Close()

	rawNotifications, err := pgx.CollectRows(rows, pgx.RowToStructByName[rawNotification])

	if err != nil {
		logrus.WithError(err).Error("failed to collect unread notifications")
		return nil, err
	}

	return lo.Map(rawNotifications, func(r rawNotification, _ int) entity.Notification {
		return r.toEntity()
	}), nil
}

func (r *NotificationRepository) MarkRead(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		UPDATE notification n
		SET status_id = ns.id,
			read_at = NOW()
		FROM notification_status ns
		WHERE n.id = $1
			AND ns.name = $2
	`

	_, err := r.GetTxManager(ctx).Exec(ctx, query, id, entity.StatusRead)

	if err != nil {

		logrus.WithField("notification_id", id.String()).
			WithError(err).
			Error("failed to mark notification read")

		return err
	}

	return nil
}

func (r *NotificationRepository) GetByID(
	ctx context.Context,
	ID uuid.UUID,
) (entity.Notification, error) {

	query, args, _ := r.Builder.
		Select(
			"n.id",
			"n.user_id",
			"n.notification_type_id",
			"nt.name as notification_type_name",
			"n.title",
			"n.body",
			"n.payload",
			"n.status_id",
			"ns.name as status_name",
			"n.created_at",
			"n.read_at",
		).
		From("notification n").
		Join("notification_type nt ON n.notification_type_id = nt.id").
		Join("notification_status ns ON n.status_id = ns.id").
		Where("n.id = ?", ID).
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)

	if err != nil {
		logrus.WithError(err).Error("failed to fetch notification")
		return entity.Notification{}, err
	}

	defer rows.Close()

	rawNotification, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[rawNotification])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Notification{}, ErrNotificationNotFound
		}
		logrus.WithError(err).Error("failed to collect notification")
		return entity.Notification{}, err
	}

	return rawNotification.toEntity(), nil
}
