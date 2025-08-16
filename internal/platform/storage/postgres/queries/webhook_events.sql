-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (
    source_id, pipeline_id, payload, original_payload, metadata, status, scheduled_at
) VALUES ($1, $2, $3, $4, COALESCE($5, '{}'::jsonb), COALESCE($6, 'pending'), $7)
RETURNING *;

-- name: GetWebhookEventByID :one
SELECT * FROM webhook_events WHERE id = $1;

-- name: ListWebhookEventsBySource :many
SELECT * FROM webhook_events
WHERE source_id = $1
ORDER BY created_at DESC;

-- name: ListWebhookEventsBySourceAndStatus :many
SELECT * FROM webhook_events
WHERE source_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: UpdateWebhookEvent :one
UPDATE webhook_events SET
    status = COALESCE($2, status),
    metadata = COALESCE($3, metadata),
    pipeline_id = COALESCE($4, pipeline_id),
    filter_results = COALESCE($5, filter_results),
    transformation_results = COALESCE($6, transformation_results),
    error_message = COALESCE($7, error_message),
    scheduled_at = COALESCE($8, scheduled_at),
    processed_at = COALESCE($9, processed_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteWebhookEvent :exec
DELETE FROM webhook_events WHERE id = $1;

-- name: ListWebhookEventsByPipeline :many
SELECT * FROM webhook_events
WHERE pipeline_id = $1
ORDER BY created_at DESC;

-- name: ListPendingWebhookEvents :many
SELECT * FROM webhook_events
WHERE status = 'pending' AND (scheduled_at IS NULL OR scheduled_at <= NOW())
ORDER BY created_at ASC
LIMIT $1;

-- name: ListFailedWebhookEvents :many
SELECT * FROM webhook_events
WHERE status = 'failed'
ORDER BY created_at DESC;

-- name: GetWebhookEventWithPipeline :one
SELECT 
    we.*,
    p.name as pipeline_name,
    s.name as source_name,
    d.name as destination_name
FROM webhook_events we
LEFT JOIN pipelines p ON we.pipeline_id = p.id
LEFT JOIN sources s ON we.source_id = s.id
LEFT JOIN destinations d ON p.destination_id = d.id
WHERE we.id = $1;

-- name: CountWebhookEventsByStatus :one
SELECT COUNT(*) FROM webhook_events WHERE status = $1;

-- name: UpdateWebhookEventStatus :one
UPDATE webhook_events SET
    status = $2,
    error_message = COALESCE($3, error_message),
    processed_at = CASE WHEN $2 IN ('delivered', 'failed', 'filtered') THEN NOW() ELSE processed_at END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetWebhookEventWithDetails :one
SELECT 
    we.*,
    p.name as pipeline_name,
    s.name as source_name,
    d.name as destination_name
FROM webhook_events we
LEFT JOIN pipelines p ON we.pipeline_id = p.id
LEFT JOIN sources s ON we.source_id = s.id
LEFT JOIN destinations d ON p.destination_id = d.id
WHERE we.id = $1;