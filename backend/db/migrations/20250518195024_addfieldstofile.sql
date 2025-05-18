-- +goose Up
-- +goose StatementBegin
ALTER TABLE files ADD COLUMN user_id INTEGER REFERENCES users(id);
ALTER TABLE files ADD COLUMN private BOOLEAN DEFAULT FALSE;
ALTER TABLE files ADD COLUMN comment TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE files DROP COLUMN user_id;
ALTER TABLE files DROP COLUMN private;
ALTER TABLE files DROP COLUMN comment;
-- +goose StatementEnd
