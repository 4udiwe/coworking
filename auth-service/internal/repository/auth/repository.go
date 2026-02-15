package auth_repository

import (
	"context"
	"errors"
	"time"

	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type AuthRepository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepository {
	return &AuthRepository{pg}
}

func (r *AuthRepository) SaveRefreshToken(
	ctx context.Context,
	userID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
) error {

	logrus.WithField("user_id", userID).Info("Saving refresh token")

	query, args, _ := r.Builder.
		Insert("refresh_tokens").
		Columns("user_id", "token_hash", "expires_at").
		Values(userID, tokenHash, expiresAt).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error("SaveRefreshToken failed")
	}
	return err
}

func (r *AuthRepository) GetUserByRefreshToken(
	ctx context.Context,
	tokenHash string,
) (uuid.UUID, error) {

	logrus.WithField("token_hash", tokenHash).Info("Looking up user by refresh token")

	query, args, _ := r.Builder.
		Select("user_id").
		From("refresh_tokens").
		Where("token_hash = ?", tokenHash).
		Where("revoked = false").
		Where("expires_at > now()").
		ToSql()

	var userID uuid.UUID
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.WithField("token_hash", tokenHash).Warn("Refresh token not found or expired")
			return uuid.Nil, ErrInvalidRefreshToken
		}
		logrus.WithError(err).WithField("token_hash", tokenHash).Error("GetUserByRefreshToken query failed")
		return uuid.Nil, err
	}

	logrus.WithField("user_id", userID).Info("Found user by refresh token")
	return userID, nil
}

func (r *AuthRepository) RevokeRefreshToken(
	ctx context.Context,
	tokenHash string,
) error {

	logrus.WithField("token_hash", tokenHash).Info("Revoking refresh token")

	query, args, _ := r.Builder.
		Update("refresh_tokens").
		Set("revoked", true).
		Where("token_hash = ?", tokenHash).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).WithField("token_hash", tokenHash).Error("RevokeRefreshToken failed")
	}
	return err
}
