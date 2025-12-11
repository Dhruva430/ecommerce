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

-- name: IncrementProductVariantStock :exec
UPDATE product_variant
SET stock = stock + $2
WHERE id = $1;

-- name: GetOrderDetailsByID :one
SELECT * FROM orders
WHERE id = $1 AND user_id = $2;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: CancelOrder :exec
UPDATE orders
SET status = 'canceled'
WHERE id = $1 AND user_id = $2;

-- name: GetOrderProducts :many
SELECT * FROM order_product
WHERE order_id = $1;