-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    provider,
    provider_id,
    role,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByProviderID :one
SELECT * FROM users
WHERE provider_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    role = COALESCE(sqlc.narg('role'), role),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET
    display_tag = sqlc.arg('display_tag'),
    colour = sqlc.arg('colour'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountBannedUsers :one
SELECT COUNT(*) FROM users WHERE banned = true;

-- name: SetUserBanned :one
UPDATE users
SET banned = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetUserMaxFileSize :one
UPDATE users
SET max_file_size_override = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
