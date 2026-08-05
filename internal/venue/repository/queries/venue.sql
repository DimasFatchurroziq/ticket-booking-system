-- name: CreateVenue :one
INSERT INTO venues (
    name,
    address,
    city,
    total_capacity
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;


-- name: GetVenueById :one
SELECT *
FROM venues
WHERE id = $1
LIMIT 1;


-- name: ListVenues :many
SELECT *
FROM venues
ORDER BY created_at DESC;


-- name: UpdateVenue :one
UPDATE venues
SET
    name = $2,
    address = $3,
    city = $4,
    total_capacity = $5
WHERE id = $1
RETURNING *;


-- name: DeleteVenue :exec
DELETE FROM venues
WHERE id = $1;


-- name: CountVenues :one
SELECT COUNT(*)
FROM venues;