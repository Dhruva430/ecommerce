-- name: CreateUser :one
INSERT INTO "user" (email, role)
VALUES ($1, $2)
RETURNING id;

-- name: CreateUserCredentials :exec
INSERT INTO user_credentials (password, first_name, last_name, phone_number, user_id)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserForLogin :one
SELECT u.id, u.email, u.role, u.verified, u.is_banned, c.password
FROM "user" u
JOIN user_credentials c ON c.user_id = u.id
WHERE u.email = $1
LIMIT 1;

-- name: AddToCart :exec
INSERT INTO cart (user_id, product_id, amount)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id)
DO UPDATE SET amount = cart.amount + excluded.amount;

-- name: UpdateCart :exec
UPDATE cart
SET amount = $3
WHERE user_id = $1
AND product_id = $2;

-- name: RemoveFromCart :exec
DELETE FROM cart WHERE user_id = $1 AND product_id = $2;

-- name: GetCart :many
SELECT c.id AS cart_id, c.amount, p.id AS product_id, p.title, p.price, p.discounted,
       p.image_url, s.shop_name
FROM cart c
JOIN product p ON p.id = c.product_id
JOIN seller s ON s.id = p.seller_id
WHERE c.user_id = $1;
