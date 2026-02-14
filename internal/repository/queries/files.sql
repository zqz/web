-- name: CreateFile :one
INSERT INTO files (
    size,
    name,
    alias,
    hash,
    slug,
    content_type,
    user_id,
    private,
    comment,
    bytes_received,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()
) RETURNING *;

-- name: GetFileByID :one
SELECT * FROM files
WHERE id = $1 LIMIT 1;

-- name: GetFileBySlug :one
SELECT * FROM files
WHERE slug = $1 LIMIT 1;

-- name: GetFileByHash :one
SELECT * FROM files
WHERE hash = $1 LIMIT 1;

-- name: ListFiles :many
SELECT * FROM files
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListFilesByUserID :many
SELECT * FROM files
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPublicFiles :many
SELECT * FROM files
WHERE private = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListFilesVisibleToUser :many
SELECT * FROM files
WHERE private = false OR user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateFile :one
UPDATE files
SET
    name = COALESCE(sqlc.narg('name'), name),
    alias = COALESCE(sqlc.narg('alias'), alias),
    slug = COALESCE(sqlc.narg('slug'), slug),
    private = COALESCE(sqlc.narg('private'), private),
    comment = COALESCE(sqlc.narg('comment'), comment),
    bytes_received = COALESCE(sqlc.narg('bytes_received'), bytes_received),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;

-- name: DeleteFilesByUserID :exec
DELETE FROM files
WHERE user_id = $1;

-- name: CountFiles :one
SELECT COUNT(*) FROM files;

-- name: CountFilesByUserID :one
SELECT COUNT(*) FROM files
WHERE user_id = $1;

-- name: TotalFileSize :one
SELECT COALESCE(SUM(size), 0)::bigint FROM files;

-- name: GetFileWithThumbnail :one
SELECT 
    f.id,
    f.size,
    f.name,
    f.alias,
    f.hash,
    f.slug,
    f.content_type,
    f.user_id,
    f.private,
    f.comment,
    f.created_at,
    f.updated_at,
    t.hash as thumbnail_hash,
    t.width as thumbnail_width,
    t.height as thumbnail_height
FROM files f
LEFT JOIN thumbnails t ON t.file_id = f.id
WHERE f.id = $1
LIMIT 1;

-- name: GetFileWithThumbnailBySlug :one
SELECT 
    f.id,
    f.size,
    f.name,
    f.alias,
    f.hash,
    f.slug,
    f.content_type,
    f.user_id,
    f.private,
    f.comment,
    f.created_at,
    f.updated_at,
    t.hash as thumbnail_hash,
    t.width as thumbnail_width,
    t.height as thumbnail_height
FROM files f
LEFT JOIN thumbnails t ON t.file_id = f.id
WHERE f.slug = $1
LIMIT 1;

-- name: GetFileWithThumbnailByHash :one
SELECT 
    f.id,
    f.size,
    f.name,
    f.alias,
    f.hash,
    f.slug,
    f.content_type,
    f.user_id,
    f.private,
    f.comment,
    f.created_at,
    f.updated_at,
    t.hash as thumbnail_hash,
    t.width as thumbnail_width,
    t.height as thumbnail_height
FROM files f
LEFT JOIN thumbnails t ON t.file_id = f.id
WHERE f.hash = $1
LIMIT 1;

-- name: ListFilesWithThumbnails :many
SELECT 
    f.id,
    f.size,
    f.name,
    f.alias,
    f.hash,
    f.slug,
    f.content_type,
    f.user_id,
    f.private,
    f.comment,
    f.created_at,
    f.updated_at,
    t.hash as thumbnail_hash,
    t.width as thumbnail_width,
    t.height as thumbnail_height
FROM files f
LEFT JOIN thumbnails t ON t.file_id = f.id
ORDER BY f.created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchFiles :many
SELECT * FROM files
WHERE (name % $1 OR alias % $1 OR COALESCE(comment, '') % $1)
   OR (POSITION(LOWER($1) IN LOWER(name)) > 0 OR POSITION(LOWER($1) IN LOWER(COALESCE(alias, ''))) > 0 OR POSITION(LOWER($1) IN LOWER(COALESCE(comment, ''))) > 0)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchPublicFiles :many
SELECT * FROM files
WHERE private = false
  AND ((name % $1 OR alias % $1 OR COALESCE(comment, '') % $1)
   OR (POSITION(LOWER($1) IN LOWER(name)) > 0 OR POSITION(LOWER($1) IN LOWER(COALESCE(alias, ''))) > 0 OR POSITION(LOWER($1) IN LOWER(COALESCE(comment, ''))) > 0))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchFilesVisibleToUser :many
SELECT * FROM files
WHERE (private = false OR user_id = $1)
  AND ((name % $2 OR alias % $2 OR COALESCE(comment, '') % $2)
   OR (POSITION(LOWER($2) IN LOWER(name)) > 0 OR POSITION(LOWER($2) IN LOWER(COALESCE(alias, ''))) > 0 OR POSITION(LOWER($2) IN LOWER(COALESCE(comment, ''))) > 0))
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
