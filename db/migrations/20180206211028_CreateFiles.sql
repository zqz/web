
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE files (
  id uuid not null default uuid_generate_v4() primary key,
  token  text not null,
  size   integer not null,
  name   text not null,
  alias  text not null,
  hash   text not null,
  content_type text not null,
  created_at timestamp without time zone,
  updated_at timestamp without time zone
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table files;

