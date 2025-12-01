-- name: CreateUser :one
INSERT INTO "user" (email, role)
VALUES ($1, $2) 
RETURNING id, email, role;

-- name: CreateAccount :exec
INSERT INTO account (account_id, provider, password, user_id)
VALUES ($1, $2, $3, $4);

-- name: GetUserByAccountID :one
SELECT u.id, u.email, u.role, a.account_id, a.provider, a.password,u.username
FROM "user" u
JOIN account a ON u.id = a.user_id
WHERE a.account_id = $1;

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