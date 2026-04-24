-- +goose Up
-- ============================================
-- Add device fingerprint column for session deduplication
-- ============================================
-- PROBLEM: 
-- When user refreshes token on same device, a NEW session is created
-- and old one marked as revoked. This creates duplicates:
-- - session-001 (revoked=true, from first login)
-- - session-002 (revoked=false, from first refresh)
-- 
-- BOTH appear on UI as separate sessions for same device!
--
-- SOLUTION:
-- Add device_fingerprint (hash of userAgent + deviceInfo)
-- This allows detecting: "Is this refresh from same device?"
--
-- If same device (fingerprint matches) → REUSE existing session
--                 (UPDATE last_used_at instead of creating new)
--
-- If new device (fingerprint different) → CREATE new session
-- ============================================
ALTER TABLE refresh_tokens
ADD COLUMN device_fingerprint VARCHAR(64);
-- Index for fast device lookup
-- Query pattern: "Find active sessions for this (user_id, device_fp)"
CREATE INDEX idx_refresh_tokens_device_fp ON refresh_tokens(user_id, device_fingerprint)
WHERE revoked = false;
-- +goose Down
DROP INDEX idx_refresh_tokens_device_fp;
ALTER TABLE refresh_tokens DROP COLUMN device_fingerprint;