-- +goose Up
ALTER TABLE files ALTER COLUMN created_at SET DEFAULT now();

-- +goose Down
ALTER TABLE files ALTER COLUMN created_at DROP DEFAULT;

