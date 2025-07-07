-- name: RegRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
);

-- name: GetUserFromRefreshToken :one
SELECT refresh_tokens.token, refresh_tokens.expires_at, refresh_tokens.revoked_at, users.id as user_id
FROM refresh_tokens
INNER JOIN users ON refresh_tokens.user_id = users.id
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;