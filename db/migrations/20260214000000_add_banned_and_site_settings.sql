-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS banned boolean NOT NULL DEFAULT false;

CREATE TABLE IF NOT EXISTS site_settings (
    key text PRIMARY KEY,
    value text NOT NULL
);

INSERT INTO site_settings (key, value) VALUES ('public_uploads_enabled', 'true')
ON CONFLICT (key) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS site_settings;
ALTER TABLE users DROP COLUMN IF EXISTS banned;
-- +goose StatementEnd
