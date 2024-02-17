
-- +goose Up
CREATE TABLE thumbnails (
  id SERIAL NOT NULL PRIMARY KEY,
  file_id INTEGER REFERENCES files(id) NOT NULL,
  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  hash TEXT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE
);

-- +goose Down
DROP TABLE thumbnails;

