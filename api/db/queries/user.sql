-- name: CreateUser :one
INSERT INTO "user" (email, role, username)
VALUES ($1, $2, $3)
RETURNING id, email, role;

-- name: GetUserByAccountID :one
SELECT u.id, u.email, u.role, a.account_id, a.provider, a.password, u.username
FROM user_view u
JOIN account a ON u.id = a.user_id
WHERE a.account_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM user_view WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, role, username
FROM user_view
WHERE id = $1;

-- name: DeleteUser :exec
UPDATE "user"
SET is_deleted = TRUE
WHERE id = $1;

-- name: GetUserAddresses :many
SELECT * FROM address
WHERE user_id = $1;

-- name: UpdateUserAddress :exec
UPDATE address
SET name = $2,
    pincode = $3,
    area = $4,
    city = $5,
    state = $6,
    country = $7,
    phone_number = $8
WHERE id = $1;

-- name: GetOrderHistory :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;