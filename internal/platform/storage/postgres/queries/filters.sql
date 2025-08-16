-- name: CreateFilter :one
INSERT INTO filters (
    pipeline_id, name, description, filter_type, mode, config, code, execution_order, is_active
) VALUES ($1, $2, COALESCE($3, ''), $4, COALESCE($5, 'nocode'), COALESCE($6, '{}'::jsonb), $7, COALESCE($8, 1), COALESCE($9, TRUE))
RETURNING *;

-- name: GetFilterByID :one
SELECT * FROM filters WHERE id = $1;

-- name: ListFiltersByPipeline :many
SELECT * FROM filters 
WHERE pipeline_id = $1 
ORDER BY execution_order ASC, created_at ASC;

-- name: ListActiveFiltersByPipeline :many
SELECT * FROM filters 
WHERE pipeline_id = $1 AND is_active = TRUE 
ORDER BY execution_order ASC;

-- name: UpdateFilter :one
UPDATE filters SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    filter_type = COALESCE($4, filter_type),
    mode = COALESCE($5, mode),
    config = COALESCE($6, config),
    code = COALESCE($7, code),
    execution_order = COALESCE($8, execution_order),
    is_active = COALESCE($9, is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteFilter :exec
DELETE FROM filters WHERE id = $1;

-- name: CountFiltersByPipeline :one
SELECT COUNT(*) FROM filters WHERE pipeline_id = $1;

-- name: GetFiltersByType :many
SELECT * FROM filters 
WHERE pipeline_id = $1 AND filter_type = $2 AND is_active = TRUE
ORDER BY execution_order ASC;

-- name: ReorderFilters :exec
UPDATE filters SET execution_order = $2, updated_at = NOW() WHERE id = $1;
