-- name: CreatePipeline :one
INSERT INTO pipelines (
    user_id, source_id, destination_id, name, description, is_active, execution_order
) VALUES ($1, $2, $3, $4, COALESCE($5, ''), COALESCE($6, TRUE), COALESCE($7, 1))
RETURNING *;

-- name: GetPipelineByID :one
SELECT * FROM pipelines WHERE id = $1;

-- name: ListPipelinesByUser :many
SELECT * FROM pipelines 
WHERE user_id = $1 
ORDER BY execution_order ASC, created_at DESC;

-- name: ListActivePipelinesBySource :many
SELECT * FROM pipelines 
WHERE source_id = $1 AND is_active = TRUE 
ORDER BY execution_order ASC;

-- name: ListPipelinesBySourceAndDestination :many
SELECT * FROM pipelines 
WHERE source_id = $1 AND destination_id = $2 
ORDER BY execution_order ASC;

-- name: UpdatePipeline :one
UPDATE pipelines SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    is_active = COALESCE($4, is_active),
    execution_order = COALESCE($5, execution_order),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePipeline :exec
DELETE FROM pipelines WHERE id = $1;

-- name: CountPipelinesByUser :one
SELECT COUNT(*) FROM pipelines WHERE user_id = $1;

-- name: GetPipelineWithDetails :one
SELECT 
    p.*,
    s.name as source_name,
    d.name as destination_name,
    d.destination_type
FROM pipelines p
JOIN sources s ON p.source_id = s.id
JOIN destinations d ON p.destination_id = d.id
WHERE p.id = $1;