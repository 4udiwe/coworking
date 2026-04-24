-- +goose Up
-- +goose StatementBegin
ALTER TABLE notification
ADD COLUMN action_url TEXT;
CREATE INDEX idx_notification_action_url ON notification (action_url)
WHERE action_url IS NOT NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_notification_action_url;
ALTER TABLE notification DROP COLUMN action_url;
-- +goose StatementEnd