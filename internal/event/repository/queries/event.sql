-- name: CreateEvent :one
INSERT INTO events (
    venue_id,
    name,
    description,
    event_start,
    event_end,
    sale_start,
    sale_end,
    status
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8
)
RETURNING *;


-- name: GetEvent :one
SELECT *
FROM events
WHERE id = $1
LIMIT 1;


-- name: ListEvents :many
SELECT *
FROM events
WHERE
    (
        sqlc.narg(filter_status)::text IS NULL
        OR status = sqlc.narg(filter_status)
    )
AND
    (
        sqlc.narg(filter_start_sale_date)::timestamptz IS NULL
        OR sale_start >= sqlc.narg(filter_start_sale_date)
    )
AND
    (
        sqlc.narg(filter_end_sale_date)::timestamptz IS NULL
        OR sale_start <= sqlc.narg(filter_end_sale_date)
    )
ORDER BY sale_start;


-- name: UpdateEvent :one
UPDATE events
SET
    venue_id = $2,
    name = $3,
    description = $4,
    event_start = $5,
    event_end = $6,
    sale_start = $7,
    sale_end = $8,
    status = $9,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;


-- name: CountEvents :one
SELECT COUNT(*)
FROM events;