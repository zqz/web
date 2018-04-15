-- +goose Up
alter table files alter column created_at set default now();

-- +goose Down
alter table files alter column created_at drop default;

