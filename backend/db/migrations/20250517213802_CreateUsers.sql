-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id SERIAL NOT NULL PRIMARY KEY,
  name text NOT NULL,
  email text NOT NULL,
  provider text NOT NULL,
  provider_id text NOT NULL,
  role text NOT NULL default 'member',
  created_at TIMESTAMP WITHOUT TIME ZONE,
  updated_at TIMESTAMP WITHOUT TIME ZONE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
