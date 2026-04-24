-- +goose Up
-- +goose StatementBegin
ALTER TABLE timer
ADD COLUMN place_id UUID;
ALTER TABLE timer
ADD COLUMN place_label VARCHAR(255);
ALTER TABLE timer
ADD COLUMN start_time TIMESTAMPTZ;
ALTER TABLE timer
ADD COLUMN end_time TIMESTAMPTZ;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE timer DROP COLUMN place_id;
ALTER TABLE timer DROP COLUMN place_label;
ALTER TABLE timer DROP COLUMN start_time;
ALTER TABLE timer DROP COLUMN end_time;
-- +goose StatementEnd