-- name: CreateRule :one
INSERT INTO rules (
    user_id, source_id, destination_id, name, version, is_active, mode, config, code
) VALUES ($1, $2, $3, $4, COALESCE($5, 1), COALESCE($6, TRUE), $7, COALESCE($8, '{}'::jsonb), $9)
RETURNING *;

-- name: GetRuleByID :one
SELECT * FROM rules WHERE id = $1;

-- name: ListActiveRulesByUser :many
SELECT * FROM rules WHERE user_id = $1 AND is_active = TRUE ORDER BY created_at DESC;

-- name: ListRulesBySourceAndDestination :many
SELECT * FROM rules WHERE source_id = $1 AND destination_id = $2 ORDER BY version DESC;

-- name: UpdateRule :one
UPDATE rules SET
    name = COALESCE($2, name),
    version = COALESCE($3, version),
    is_active = COALESCE($4, is_active),
    mode = COALESCE($5, mode),
    config = COALESCE($6, config),
    code = COALESCE($7, code),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteRule :exec
DELETE FROM rules WHERE id = $1;
