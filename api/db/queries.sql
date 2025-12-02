-- name: CreateUser :one
INSERT INTO "user" (email, role)
VALUES ($1, $2) 
RETURNING id, email, role;

-- name: CreateAccount :exec
INSERT INTO account (account_id, provider, password, user_id)
VALUES ($1, $2, $3, $4);

-- name: GetUserByAccountID :one
SELECT u.id, u.email, u.role, a.account_id, a.provider, a.password,u.username
FROM "user_view" u
JOIN account a ON u.id = a.user_id
WHERE a.account_id = $1;

-- name: GetUserByID :one
SELECT id, email, role, username
FROM "user_view"
WHERE id = $1;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_token (id,token, user_id, ip_address, expires_at)
VALUES ($1, $2, $3, $4,$5);

-- name: DeleteRefreshTokensByID :exec
DELETE FROM refresh_token
WHERE id = $1;

-- name: UpdateRefreshTokenRevoked :exec
UPDATE refresh_token
SET revoked = $2 , last_used = $3
WHERE id = $1;

-- name: GetRefreshToken :one
SELECT id, token, user_id, revoked, ip_address, created_at, expires_at
FROM refresh_token
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


-- name: DeleteUser :exec
UPDATE "user"
SET is_deleted = TRUE
WHERE id = $1;