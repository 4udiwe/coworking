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
			"device_fingerprint",
			"token_hash",
			"expires_at",
		).
		Values(
			session.ID,
			session.UserID,
			session.UserAgent,
			session.IPAddress,
			session.DeviceName,
			session.DeviceFingerprint,
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
			"device_fingerprint",
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
		&s.DeviceFingerprint,
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
			"device_fingerprint",
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
			&s.DeviceFingerprint,
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
		Where("revoked = false").  // Only active sessions
		OrderBy("created_at ASC"). // Order by oldest first
		Limit(1).                  // Get only the first one
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

// GetSessionByDeviceFingerprint finds an active session by device fingerprint and user ID
//
// PURPOSE:
// Implements session reuse detection. When a user issues a refresh token request,
// we check if they already have an active session on the same device (identified by fingerprint).
//
// HOW IT WORKS:
// 1. Find the most recent non-revoked, non-expired session for the user
// 2. Filter by device_fingerprint to ensure it's the same physical device
// 3. Return the session so it can be reused instead of creating a duplicate
//
// IMPORTANT:
// - Returns "no rows in result set" error if no matching session found (caller must handle)
// - Only returns active sessions (revoked=false AND expires_at > now())
// - Handles NULL device_fingerprint (old sessions without fingerprint won't match)
//
// EXAMPLE:
// User refreshes token on iPhone at 15:00
// - iPhone session exists from 10:00 (device_fingerprint = "abc123")
// - Same fingerprint "abc123" detected
// - Return existing session → Reuse it (update last_used_at)
// → UI shows 1 session, not 2 duplicates
//
// Different User refreshes on Android at 15:15
// - iPhone session exists (fingerprint = "abc123")
// - Android fingerprint = "xyz789" (different)
// - No matching session found
// - Create new session for Android
// → User has 2 active sessions (iPhone + Android)
func (r *AuthRepository) GetSessionByDeviceFingerprint(
	ctx context.Context,
	userID uuid.UUID,
	deviceFingerprint string,
) (entity.Session, error) {

	query, args, _ := r.Builder.
		Select(
			"id",
			"user_id",
			"user_agent",
			"ip_address",
			"device_name",
			"device_fingerprint",
			"expires_at",
			"last_used_at",
			"revoked",
			"created_at",
		).
		From("refresh_tokens").
		Where("user_id = ?", userID).
		Where("device_fingerprint = ?", deviceFingerprint).
		Where("revoked = false").    // Active only
		Where("expires_at > now()"). // Not expired
		OrderBy("created_at DESC").  // Get most recent
		Limit(1).
		ToSql()

	var s entity.Session
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&s.ID,
		&s.UserID,
		&s.UserAgent,
		&s.IPAddress,
		&s.DeviceName,
		&s.DeviceFingerprint,
		&s.ExpiresAt,
		&s.LastUsedAt,
		&s.Revoked,
		&s.CreatedAt,
	)

	return s, err
}

// UpdateSessionRefresh reuses an existing session by updating its token hash and last_used_at
//
// PURPOSE:
// When a user refreshes their token on the same device (detected via device fingerprint),
// instead of creating a new session (which looks like a duplicate), we update the existing
// session to mark it as recently used.
//
// HOW IT WORKS:
// 1. Update token_hash to the new token value (old token becomes invalid)
// 2. Update expires_at to extend expiration
// 3. Update last_used_at to now() (shows when last used)
// 4. Keep id and created_at unchanged (maintains session identity and original creation time)
//
// IMPORTANT:
// - DOES NOT revoke the old session (it gets replaced immediately)
// - Preserves createdAt (timestamp of original login, not this refresh)
// - Updates lastUsedAt to current time (when this refresh happened)
// - Result: UI shows 1 session with correct timestamps, no duplicates
//
// EXAMPLE - Token Refresh on Same Device:
// Before call:
//   Session ID: abc-123
//   Created: 2026-03-30 10:00
//   LastUsed: 2026-03-30 10:00  ← old
//   Expires: 2026-04-13 10:00   ← old
//   TokenHash: hash_of_old_token
//
// After call to UpdateSessionRefresh():
//   Session ID: abc-123         ← SAME (no duplicate)
//   Created: 2026-03-30 10:00   ← UNCHANGED (correct createdAt)
//   LastUsed: 2026-03-30 16:00  ← UPDATED (when refresh happened)
//   Expires: 2026-04-13 16:00   ← EXTENDED
//   TokenHash: hash_of_new_token
//
// Result: User sees single session on UI with all correct timestamps
func (r *AuthRepository) UpdateSessionRefresh(
	ctx context.Context,
	sessionID uuid.UUID,
	newTokenHash string,
	newExpiresAt time.Time,
) error {

	query, args, _ := r.Builder.
		Update("refresh_tokens").
		Set("token_hash", newTokenHash).
		Set("expires_at", newExpiresAt).
		Set("last_used_at", time.Now()).
		Where("id = ?", sessionID).
		ToSql()

	_, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	return err
}

// DeleteOldRevokedSessions удаляет revoked сессии старше заданного количества дней
//
// Cleanup старых revoked сессий из базы данных.
// Вызывается периодически через Kafka события от scheduler-service.
//
// 1. Получает retentionDays (например, 10)
// 2. Удаляет все сессии WHERE revoked = true AND created_at < now() - INTERVAL 'retentionDays days'
// 3. Возвращает количество удалённых строк
func (r *AuthRepository) DeleteOldRevokedSessions(
	ctx context.Context,
	retentionDays int,
) (int64, error) {
	query, args, _ := r.Builder.
		Delete("refresh_tokens").
		Where("revoked = true").
		Where("created_at < now() - INTERVAL '1 day' * ?", retentionDays).
		ToSql()

	result, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected := result.RowsAffected()
	return rowsAffected, nil
}
