-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN display_tag VARCHAR(3) DEFAULT NULL,
ADD COLUMN colour VARCHAR(7) DEFAULT NULL;

-- Optional: constrain colour to hex format (#RRGGBB)
ALTER TABLE users ADD CONSTRAINT users_colour_hex_check
CHECK (colour IS NULL OR colour ~ '^#[0-9A-Fa-f]{6}$');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_colour_hex_check;
ALTER TABLE users DROP COLUMN IF EXISTS display_tag;
ALTER TABLE users DROP COLUMN IF EXISTS colour;
-- +goose StatementEnd
