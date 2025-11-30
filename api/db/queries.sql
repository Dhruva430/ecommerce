-- name: CreateUser :one
INSERT INTO "user" (email, role)
VALUES ($1, $2) 
RETURNING id, email, role;

-- name: CreateAccount :exec
INSERT INTO account (account_id, provider, password, user_id)
VALUES ($1, $2, $3, $4);

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
