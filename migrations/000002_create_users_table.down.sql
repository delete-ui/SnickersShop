-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;

-- Удаляем функцию, если она осталась
DROP FUNCTION IF EXISTS update_users_updated_at();
-- +goose StatementEnd