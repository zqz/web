-- +goose Up
-- +goose StatementBegin
ALTER TABLE files ADD COLUMN bytes_received INTEGER NOT NULL DEFAULT 0;

-- For existing files, set bytes_received equal to size (assume they're complete)
UPDATE files SET bytes_received = size WHERE bytes_received = 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE files DROP COLUMN bytes_received;
-- +goose StatementEnd
