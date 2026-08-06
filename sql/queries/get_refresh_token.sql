-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens WHERE token = $1 AND revoked_at IS NOT NULL AND expires_at >= NOW();