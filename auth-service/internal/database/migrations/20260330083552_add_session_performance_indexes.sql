-- +goose Up
-- Index 1: Simple index on expires_at
-- Used for filtering expired sessions
ALTER TABLE refresh_tokens
ADD INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
-- Index 2: Composite index (user_id, expires_at DESC)
-- Optimal for query: WHERE user_id=X AND expires_at > now()
ALTER TABLE refresh_tokens
ADD INDEX idx_refresh_tokens_user_expires ON refresh_tokens(user_id, expires_at DESC);
-- +goose Down
DROP INDEX idx_refresh_tokens_expires_at ON refresh_tokens;
DROP INDEX idx_refresh_tokens_user_expires ON refresh_tokens;