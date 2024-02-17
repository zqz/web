
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE files (
  id SERIAL NOT NULL PRIMARY KEY,
  size   INTEGER NOT NULL,
  name   TEXT NOT NULL,
  alias  TEXT NOT NULL,
  hash   TEXT NOT NULL,
  slug    TEXT NOT NULL,
  content_type TEXT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE,
  updated_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX idx_slug_on_files ON files (slug);
CREATE INDEX idx_hash_on_files ON files (hash);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE files;

