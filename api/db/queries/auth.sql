-- name: CreateAccount :exec
INSERT INTO account (account_id, provider, password, user_id)
VALUES ($1, $2, $3, $4);

-- name: CreateRefreshToken :exec
INSERT INTO refresh_token (id, token, user_id, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateRefreshTokenRevoked :exec
UPDATE refresh_token
SET revoked = $2, last_used = $3
WHERE id = $1;

-- name: GetRefreshToken :one
SELECT id, token, user_id, revoked, ip_address, created_at, expires_at
FROM refresh_token
WHERE id = $1;

-- name: DeleteRefreshTokensByID :exec
DELETE FROM refresh_token WHERE id = $1;

-- name: CreateBuyer :one
INSERT INTO buyer (user_id)
VALUES ($1)
RETURNING *;