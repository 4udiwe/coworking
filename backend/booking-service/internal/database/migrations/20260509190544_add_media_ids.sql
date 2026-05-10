-- +goose Up
ALTER TABLE coworking
ADD COLUMN media_ids text[] NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE coworking
DROP COLUMN media_ids;
