-- +goose Up
ALTER TABLE coworking_layout
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT false;
-- добавление поля is_active к лейауту для отслеживания активного 
CREATE UNIQUE INDEX uniq_active_layout ON coworking_layout(coworking_id)
WHERE is_active = true;
-- индекс гарантирует что у одного коворкинга только один активный лейаут
-- +goose Down
ALTER TABLE coworking_layout DROP COLUMN is_active;
DROP INDEX IF EXISTS uniq_active_layout;