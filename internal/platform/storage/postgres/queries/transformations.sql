-- name: CreateTransformation :one
INSERT INTO transformations (
    pipeline_id, name, description, transformation_type, mode, config, code, execution_order, is_active
) VALUES ($1, $2, COALESCE($3, ''), $4, COALESCE($5, 'nocode'), COALESCE($6, '{}'::jsonb), $7, COALESCE($8, 1), COALESCE($9, TRUE))
RETURNING *;

-- name: GetTransformationByID :one
SELECT * FROM transformations WHERE id = $1;

-- name: ListTransformationsByPipeline :many
SELECT * FROM transformations 
WHERE pipeline_id = $1 
ORDER BY execution_order ASC, created_at ASC;

-- name: ListActiveTransformationsByPipeline :many
SELECT * FROM transformations 
WHERE pipeline_id = $1 AND is_active = TRUE 
ORDER BY execution_order ASC;

-- name: UpdateTransformation :one
UPDATE transformations SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    transformation_type = COALESCE($4, transformation_type),
    mode = COALESCE($5, mode),
    config = COALESCE($6, config),
    code = COALESCE($7, code),
    execution_order = COALESCE($8, execution_order),
    is_active = COALESCE($9, is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTransformation :exec
DELETE FROM transformations WHERE id = $1;

-- name: CountTransformationsByPipeline :one
SELECT COUNT(*) FROM transformations WHERE pipeline_id = $1;

-- name: GetTransformationsByType :many
SELECT * FROM transformations 
WHERE pipeline_id = $1 AND transformation_type = $2 AND is_active = TRUE
ORDER BY execution_order ASC;

-- name: ReorderTransformations :exec
UPDATE transformations SET execution_order = $2, updated_at = NOW() WHERE id = $1;

-- name: GetHeaderTransformations :many
SELECT * FROM transformations 
WHERE pipeline_id = $1 AND transformation_type IN ('header_add', 'header_remove', 'header_modify') AND is_active = TRUE
ORDER BY execution_order ASC;

-- name: GetBodyTransformations :many
SELECT * FROM transformations 
WHERE pipeline_id = $1 AND transformation_type IN ('body_add', 'body_remove', 'body_modify') AND is_active = TRUE
ORDER BY execution_order ASC;
