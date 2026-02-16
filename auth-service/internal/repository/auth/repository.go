package auth_repository

import (
	"context"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/4udiwe/coworking/auth-service/internal/entity"
	"github.com/google/uuid"
)

type AuthRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepository {
	return &AuthRepository{pg}
}

func (r *AuthRepository) CreateSession(
	ctx context.Context,
	session entity.Session,
	tokenHash string,
) error {

	query, args, _ := r.Builder.
		Insert("refresh_tokens").
		Columns(
			"id",
			"user_id",
			"user_agent",
			"ip_address",
			"device_name",
			"token_hash",
			"expires_at",
		).
		Values(
			session.ID,
			session.UserID,
			session.UserAgent,
			session.IPAddress,
			session.DeviceName,
			tokenHash,
			session.ExpiresAt,
		).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	return err
}

func (r *AuthRepository) GetSessionByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Session, error) {

	query, args, _ := r.Builder.
		Select(
			"id",
			"user_id",
			"user_agent",
			"ip_address",
			"device_name",
			"expires_at",
			"last_used_at",
			"revoked",
			"created_at",
		).
		From("refresh_tokens").
		Where("id = ?", id).
		ToSql()

	var s entity.Session
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&s.ID,
		&s.UserID,
		&s.UserAgent,
		&s.IPAddress,
		&s.DeviceName,
		&s.ExpiresAt,
		&s.LastUsedAt,
		&s.Revoked,
		&s.CreatedAt,
	)

	return s, err
}

func (r *AuthRepository) UpdateLastUsedAt(
	ctx context.Context,
	id uuid.UUID,
) error {

	query, args, _ := r.Builder.
		Update("refresh_tokens").
		Set("last_used_at", time.Now()).
		Where("id = ?", id).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	return err
}

func (r *AuthRepository) RevokeSession(
	ctx context.Context,
	id uuid.UUID,
) error {

	query, args, _ := r.Builder.
		Update("refresh_tokens").
		Set("revoked", true).
		Where("id = ?", id).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	return err
}

func (r *AuthRepository) GetUserSessions(
	ctx context.Context,
	userID uuid.UUID,
	onlyActive bool,
) ([]entity.Session, error) {

	builder := r.Builder.
		Select(
			"id",
			"user_agent",
			"ip_address",
			"device_name",
			"expires_at",
			"last_used_at",
			"revoked",
			"created_at",
		).
		From("refresh_tokens").
		Where("user_id = ?", userID)

	if onlyActive {
		builder = builder.
			Where("revoked = false").
			Where("expires_at > now()")
	}

	query, args, _ := builder.
		OrderBy("created_at DESC").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []entity.Session

	for rows.Next() {
		var s entity.Session
		err := rows.Scan(
			&s.ID,
			&s.UserAgent,
			&s.IPAddress,
			&s.DeviceName,
			&s.ExpiresAt,
			&s.LastUsedAt,
			&s.Revoked,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}
