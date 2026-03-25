-- +goose Up
CREATE INDEX idx_users_name ON users(first_name, last_name);
CREATE INDEX idx_roles_code ON roles(code);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
-- +goose Down
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_roles_code;
DROP INDEX IF EXISTS idx_user_roles_user_id;