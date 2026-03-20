-- +goose Up
ALTER TABLE users
ADD COLUMN first_name TEXT;
ALTER TABLE users
ADD COLUMN last_name TEXT;
-- +goose Down
ALTER TABLE users DROP COLUMN first_name;
ALTER TABLE users DROP COLUMN last_name;