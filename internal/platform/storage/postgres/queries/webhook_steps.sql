-- name: CreateWebhookStep :one
INSERT INTO webhook_steps (
    webhook_event_id, pipeline_id, step_type, step_name, step_id, execution_order, 
    status, input_data, output_data, error_message, duration_ms, started_at, completed_at
) VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, 'pending'), COALESCE($8, '{}'::jsonb), 
          COALESCE($9, '{}'::jsonb), $10, $11, COALESCE($12, NOW()), $13)
RETURNING *;

-- name: GetWebhookStepByID :one
SELECT * FROM webhook_steps WHERE id = $1;

-- name: ListWebhookStepsByEvent :many
SELECT * FROM webhook_steps 
WHERE webhook_event_id = $1 
ORDER BY execution_order ASC, started_at ASC;

-- name: ListWebhookStepsByEventAndType :many
SELECT * FROM webhook_steps 
WHERE webhook_event_id = $1 AND step_type = $2
ORDER BY execution_order ASC, started_at ASC;

-- name: UpdateWebhookStep :one
UPDATE webhook_steps SET
    status = COALESCE($2, status),
    output_data = COALESCE($3, output_data),
    error_message = COALESCE($4, error_message),
    duration_ms = COALESCE($5, duration_ms),
    completed_at = COALESCE($6, completed_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateWebhookStepStatus :one
UPDATE webhook_steps SET
    status = $2,
    error_message = COALESCE($3, error_message),
    completed_at = CASE 
        WHEN $2 IN ('success', 'failed', 'skipped') AND completed_at IS NULL 
        THEN NOW() 
        ELSE completed_at 
    END,
    duration_ms = CASE 
        WHEN $2 IN ('success', 'failed', 'skipped') AND duration_ms IS NULL
        THEN EXTRACT(EPOCH FROM (NOW() - started_at)) * 1000
        ELSE duration_ms
    END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteWebhookStep :exec
DELETE FROM webhook_steps WHERE id = $1;

-- name: GetWebhookTraceComplete :many
SELECT 
    ws.*,
    we.source_id,
    s.name as source_name,
    p.name as pipeline_name,
    d.name as destination_name
FROM webhook_steps ws
JOIN webhook_events we ON ws.webhook_event_id = we.id
LEFT JOIN sources s ON we.source_id = s.id
LEFT JOIN pipelines p ON ws.pipeline_id = p.id
LEFT JOIN destinations d ON p.destination_id = d.id
WHERE ws.webhook_event_id = $1
ORDER BY ws.execution_order ASC, ws.started_at ASC;

-- name: GetFailedWebhookSteps :many
SELECT 
    ws.*,
    we.source_id,
    s.name as source_name,
    p.name as pipeline_name
FROM webhook_steps ws
JOIN webhook_events we ON ws.webhook_event_id = we.id
LEFT JOIN sources s ON we.source_id = s.id
LEFT JOIN pipelines p ON ws.pipeline_id = p.id
WHERE ws.status = 'failed'
ORDER BY ws.started_at DESC
LIMIT $1;
