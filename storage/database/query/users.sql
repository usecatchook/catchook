-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, avatar_url, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = true;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = true;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(sqlc.narg('name'), name),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at = NOW()
WHERE id = $1 AND is_active = true
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = $1;
