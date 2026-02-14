-- name: CreateThumbnail :one
INSERT INTO thumbnails (
    file_id,
    hash,
    width,
    height
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetThumbnailByFileID :one
SELECT * FROM thumbnails
WHERE file_id = $1 LIMIT 1;

-- name: GetThumbnailsByFileID :many
SELECT * FROM thumbnails
WHERE file_id = $1;

-- name: DeleteThumbnailsByFileID :exec
DELETE FROM thumbnails
WHERE file_id = $1;

-- name: DeleteThumbnail :exec
DELETE FROM thumbnails
WHERE id = $1;

-- name: UpdateThumbnail :one
UPDATE thumbnails
SET
    hash = COALESCE(sqlc.narg('hash'), hash),
    width = COALESCE(sqlc.narg('width'), width),
    height = COALESCE(sqlc.narg('height'), height)
WHERE id = sqlc.arg('id')
RETURNING *;
