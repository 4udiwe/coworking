-- +goose Up
-- ============================================
-- Enable automatic daily cleanup of old revoked sessions
-- ============================================
-- HOW IT WORKS:
-- 1. pg_cron extension (PostgreSQL scheduler) runs SQL periodically
-- 2. Every day at 2:00 AM: DELETE revoked sessions > 10 days old
-- 3. New revoked sessions kept for 10 days (audit window)
-- 4. No new code to maintain, no extra service to deploy
-- ============================================
-- Step 1: Ensure pg_cron extension is installed
CREATE EXTENSION IF NOT EXISTS pg_cron;
-- Step 2: Create procedure for the cleanup logic
-- A procedure is like a function that can be scheduled
CREATE OR REPLACE PROCEDURE cleanup_old_revoked_sessions() LANGUAGE SQL AS $$
DELETE FROM refresh_tokens
WHERE revoked = true
    AND created_at < now() - INTERVAL '10 days';
$$;
-- Step 3: Schedule the job
-- Cron format: minute hour day_of_month month day_of_week
-- '0 2 * * *' = at 2:00 AM every day
-- More info: https://github.com/citusdata/pg_cron
SELECT cron.schedule(
        'cleanup-revoked-sessions',
        '0 2 * * *',
        'CALL cleanup_old_revoked_sessions()'
    );
-- +goose Down
-- Unschedule the job
SELECT cron.unschedule('cleanup-revoked-sessions');
DROP PROCEDURE cleanup_old_revoked_sessions();