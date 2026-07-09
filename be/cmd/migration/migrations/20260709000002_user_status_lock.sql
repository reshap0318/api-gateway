-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active' AFTER password,
    ADD COLUMN failed_login_attempts INT UNSIGNED NOT NULL DEFAULT 0 AFTER status,
    ADD COLUMN locked_until DATETIME(3) NULL AFTER failed_login_attempts;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN status,
    DROP COLUMN failed_login_attempts,
    DROP COLUMN locked_until;
-- +goose StatementEnd
