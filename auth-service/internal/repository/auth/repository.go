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
			&s.UserID,
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

// DeleteOldestSessionByUser revokes the oldest active (non-revoked) session for a user
// 
// PURPOSE:
// Enforces MAX_SESSIONS_PER_USER limit. When a user exceeds the limit,
// this method finds and revokes their oldest session to make room for a new one.
//
// HOW IT WORKS:
// 1. Find the oldest non-revoked session (WHERE revoked = false)
// 2. Sort by created_at ASC (oldest first)
// 3. Revoke that session (SET revoked = true)
//
// EXAMPLE:
// User has 5 active sessions (at limit):
//   [NEW]    iPhone 2026-03-30 10:00  ← oldest
//   Web      2026-03-25 14:30
//   iPad     2026-03-28 09:15
//   MacBook  2026-03-28 16:45
//   Android  2026-03-29 12:00        ← newest
//
// When user logins on 6th device (Windows):
//   → This method finds iPhone session
//   → Revokes it (created_at = 2026-03-30)
//   → Now there are 4 active + 1 new = 5 total
//
// IMPORTANT: Returns nil if no active sessions found (no error!)
func (r *AuthRepository) DeleteOldestSessionByUser(
	ctx context.Context,
	userID uuid.UUID,
) error {
	// Step 1: Find the oldest non-revoked session
	query, args, _ := r.Builder.
		Select("id").
		From("refresh_tokens").
		Where("user_id = ?", userID).
		Where("revoked = false").           // Only active sessions
		OrderBy("created_at ASC").          // Order by oldest first
		Limit(1).                            // Get only the first one
		ToSql()

	// Step 2: Execute query to get session ID
	var oldestSessionID uuid.UUID
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&oldestSessionID)
	if err != nil {
		// If no rows found = no active sessions = OK (not an error)
		if err.Error() == "no rows in result set" {
			return nil
		}
		return err
	}

	// Step 3: Revoke the oldest session
	// This just calls RevokeSession() which sets revoked = true
	return r.RevokeSession(ctx, oldestSessionID)
}
