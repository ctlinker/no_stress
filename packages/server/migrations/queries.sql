-- name: CreateUser :execresult
INSERT INTO user (
  name, email, password_hash, updated_at
) VALUES (
  ?, ?, ?, NOW()
);

-- name: GetUser :one
SELECT * FROM user
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: ShowSoftDeletedUsers :many
SELECT * FROM user
WHERE deleted_at IS NOT NULL;

-- name: HardDeleteUser :exec
DELETE FROM user
WHERE id = ?;

-- name: SoftDeleteUser :exec
UPDATE user
SET deleted_at = NOW()
WHERE id = ?;

-- name: GetUserByMail :one
SELECT * FROM user
WHERE email = ?
LIMIT 1;

-- name: CreateSession :execresult
INSERT INTO session (
  user_id, token_hash, updated_at
) VALUES (
  ?, ?, ?
);

-- name: GetSessionByTokenHash :one
SELECT * FROM session
WHERE token_hash = ? AND deleted_at IS NULL
LIMIT 1;

-- name: SoftDeleteSessionByTokenHash :exec
UPDATE session
SET deleted_at = NOW()
WHERE token_hash = ?;