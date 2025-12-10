-- name: CreateOrder :one
INSERT INTO orders (user_id, address_id, total_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateOrderProduct :one
INSERT INTO order_product (order_id, product_id, variant_id, amount)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DecrementProductVariantStock :exec
UPDATE product_variant
SET stock = stock - $2
WHERE id = $1 AND stock >= $2;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 AND user_id = $2;