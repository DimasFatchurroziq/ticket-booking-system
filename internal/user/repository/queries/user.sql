-- ==========================================
-- Create
-- ==========================================

-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    full_name,
    phone_number
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- ==========================================
-- Read
-- ==========================================

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- ==========================================
-- Update
-- ==========================================

-- name: UpdateUser :one
UPDATE users
SET
    full_name = $2,
    phone_number = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET
    password_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserEmail :one
UPDATE users
SET
    email = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- ==========================================
-- Delete
-- ==========================================

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- ==========================================
-- Utility
-- ==========================================

-- name: CountUsers :one
SELECT COUNT(*)
FROM users;

-- name: ExistsUserByEmail :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE email = $1
);