// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT refresh_tokens.token, refresh_tokens.expires_at, refresh_tokens.revoked_at, users.id as user_id
FROM refresh_tokens
INNER JOIN users ON refresh_tokens.user_id = users.id
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL
`

type GetUserFromRefreshTokenRow struct {
	Token     string
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	UserID    uuid.UUID
}

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (GetUserFromRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i GetUserFromRefreshTokenRow
	err := row.Scan(
		&i.Token,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const regRefreshToken = `-- name: RegRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
)
`

type RegRefreshTokenParams struct {
	Token  string
	UserID uuid.UUID
}

func (q *Queries) RegRefreshToken(ctx context.Context, arg RegRefreshTokenParams) error {
	_, err := q.db.ExecContext(ctx, regRefreshToken, arg.Token, arg.UserID)
	return err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshToken, token)
	return err
}
