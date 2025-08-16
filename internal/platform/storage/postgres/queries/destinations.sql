-- name: CreateDestination :one
INSERT INTO destinations (
    user_id, name, description, destination_type, config, is_active, delay_seconds, retry_attempts
) VALUES ($1, $2, $3, $4, COALESCE($5, '{}'::jsonb), COALESCE($6, TRUE), COALESCE($7, 0), COALESCE($8, 0))
RETURNING *;

-- name: GetDestinationByID :one
SELECT * FROM destinations WHERE id = $1;

-- name: GetDestinationByName :one
SELECT * FROM destinations WHERE name = $1;

-- name: ListDestinations :many
SELECT name, description, destination_type, is_active, created_at, updated_at FROM destinations
WHERE 
    ($1 = '' OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
    AND ($2 = '' OR destination_type = $2::destination_type)
    AND (NOT $3 OR is_active = $4)
ORDER BY 
    CASE 
        WHEN $5 = 'name' AND $6 = 'asc' THEN name
    END ASC,
    CASE 
        WHEN $5 = 'name' AND $6 = 'desc' THEN name
    END DESC,
    CASE 
        WHEN $5 = 'created_at' AND $6 = 'asc' THEN created_at
    END ASC,
    CASE 
        WHEN $5 = 'updated_at' AND $6 = 'asc' THEN updated_at
    END ASC,
    CASE 
        WHEN $5 = 'updated_at' AND $6 = 'desc' THEN updated_at
    END DESC,
    CASE 
        WHEN $5 = 'is_active' AND $6 = 'asc' THEN is_active
    END ASC,
    CASE 
        WHEN $5 = 'is_active' AND $6 = 'desc' THEN is_active
    END DESC,
    CASE 
        WHEN $5 = 'created_at' AND $6 = 'desc' OR $5 = '' OR $5 IS NULL THEN created_at
    END DESC
LIMIT $7 OFFSET $8;

-- name: CountDestinations :one
SELECT COUNT(*) FROM destinations
WHERE 
    ($1 = '' OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
    AND ($2 = '' OR destination_type = $2::destination_type)
    AND (NOT $3 OR is_active = $4);

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
