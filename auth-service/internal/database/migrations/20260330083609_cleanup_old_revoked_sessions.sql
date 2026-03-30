-- +goose Up
-- ============================================
-- One-time cleanup of ALL revoked sessions
-- ============================================
DELETE FROM refresh_tokens
WHERE revoked = true;
-- +goose Down