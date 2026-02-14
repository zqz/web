-- +goose Up
-- +goose StatementBegin
-- Default max file size in bytes (100 MB). Stored in site_settings.
INSERT INTO site_settings (key, value) VALUES ('default_max_file_size', '104857600')
ON CONFLICT (key) DO NOTHING;

-- Per-user override (null = use site default). Only settable by admins.
ALTER TABLE users ADD COLUMN IF NOT EXISTS max_file_size_override BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM site_settings WHERE key = 'default_max_file_size';
ALTER TABLE users DROP COLUMN IF EXISTS max_file_size_override;
-- +goose StatementEnd
