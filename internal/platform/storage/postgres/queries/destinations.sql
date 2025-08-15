-- name: CreateDestination :one
INSERT INTO destinations (
    user_id, name, description, destination_type, config, is_active, delay_seconds, retry_attempts
) VALUES ($1, $2, $3, $4, COALESCE($5, '{}'::jsonb), COALESCE($6, TRUE), COALESCE($7, 0), COALESCE($8, 0))
RETURNING *;

-- name: GetDestinationByID :one
SELECT * FROM destinations WHERE id = $1;

-- name: ListActiveDestinationsByUser :many
SELECT * FROM destinations WHERE user_id = $1 AND is_active = TRUE ORDER BY created_at DESC;

-- name: UpdateDestination :one
UPDATE destinations SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    destination_type = COALESCE($4, destination_type),
    config = COALESCE($5, config),
    is_active = COALESCE($6, is_active),
    delay_seconds = COALESCE($7, delay_seconds),
    retry_attempts = COALESCE($8, retry_attempts),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDestination :exec
DELETE FROM destinations WHERE id = $1;
