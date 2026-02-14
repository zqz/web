-- name: GetSiteSetting :one
SELECT value FROM site_settings WHERE key = $1;

-- name: SetSiteSetting :exec
INSERT INTO site_settings (key, value) VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE SET value = $2;
