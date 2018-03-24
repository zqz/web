
-- +goose Up
CREATE TABLE thumbnails (
  id serial not null primary key,
  file_id integer references files(id),
  width integer not null,
  height integer not null,
  hash text not null,
  created_at timestamp without time zone
);

-- +goose Down
drop table thumbnails;

