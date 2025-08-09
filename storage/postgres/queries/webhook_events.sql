-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (
    source_id, payload, metadata, applied_rule_version_id, status, scheduled_at
) VALUES ($1, $2, COALESCE($3, '{}'::jsonb), $4, $5, $6)
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
    applied_rule_version_id = COALESCE($4, applied_rule_version_id),
    scheduled_at = COALESCE($5, scheduled_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteWebhookEvent :exec
DELETE FROM webhook_events WHERE id = $1;