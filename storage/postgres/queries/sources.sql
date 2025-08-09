-- name: CreateSource :one
INSERT INTO sources (
    name, user_id, description, protocol, auth_type, auth_config, is_active
) VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, TRUE))
RETURNING *;

-- name: GetSourceByID :one
SELECT * FROM sources WHERE id = $1;

-- name: ListSourcesByUser :many
SELECT * FROM sources WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateSource :one
UPDATE sources SET
   name = COALESCE($2, name),
   description = COALESCE($3, description),
   protocol = COALESCE($4, protocol),
   auth_type = COALESCE($5, auth_type),
   auth_config = COALESCE($6, auth_config),
   is_active = COALESCE($7, is_active),
   updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteSource :exec
DELETE FROM sources WHERE id = $1;

-- name: GetSourceByName :one
SELECT * FROM sources where name = $1;