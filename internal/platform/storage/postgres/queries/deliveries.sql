-- name: CreateDelivery :one
INSERT INTO deliveries (
    webhook_event_id, destination_id, status, response_code, attempt, last_error, scheduled_at
) VALUES ($1, $2, $3, $4, COALESCE($5, 0), $6, $7)
RETURNING *;

-- name: GetDeliveryByID :one
SELECT * FROM deliveries WHERE id = $1;

-- name: ListDeliveriesByWebhookEvent :many
SELECT * FROM deliveries WHERE webhook_event_id = $1 ORDER BY created_at DESC;

-- name: UpdateDelivery :one
UPDATE deliveries SET
    status = COALESCE($2, status),
    response_code = COALESCE($3, response_code),
    attempt = COALESCE($4, attempt),
    last_error = COALESCE($5, last_error),
    scheduled_at = COALESCE($6, scheduled_at),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDelivery :exec
DELETE FROM deliveries WHERE id = $1;
