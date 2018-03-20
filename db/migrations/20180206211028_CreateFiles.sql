
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE files (
  id serial not null primary key,
  size   integer not null,
  name   text not null,
  alias  text not null,
  hash   text not null,
  slug    text not null,
  content_type text not null,
  created_at timestamp without time zone,
  updated_at timestamp without time zone
);

create index idx_slug_on_files on files (slug);
create index idx_hash_on_files on files (hash);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table files;

