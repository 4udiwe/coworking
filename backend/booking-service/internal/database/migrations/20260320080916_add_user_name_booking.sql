-- +goose Up
ALTER TABLE booking
ADD COLUMN user_name TEXT DEFAULT 'Jhon Doe';
-- +goose Down
ALTER TABLE booking DROP COLUMN user_name;