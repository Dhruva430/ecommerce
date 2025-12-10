-- name: GetAllProducts :many
SELECT * FROM product
WHERE (sqlc.narg(category_id)::text IS NULL OR category_id = sqlc.arg(category_id))
  AND (sqlc.narg(seller_id)::bigint IS NULL OR seller_id = sqlc.arg(seller_id))
  AND (sqlc.narg(is_active)::boolean IS NULL OR is_active = sqlc.arg(is_active))
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: CountProducts :one
SELECT COUNT(*) FROM product
WHERE (sqlc.narg(category_id)::text IS NULL OR category_id = sqlc.arg(category_id))
  AND (sqlc.narg(seller_id)::bigint IS NULL OR seller_id = sqlc.arg(seller_id))
  AND (sqlc.narg(is_active)::boolean IS NULL OR is_active = sqlc.arg(is_active));

-- name: GetProductByID :one
SELECT * FROM product
WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO product (title, description, category_id, seller_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateProductVariant :one
INSERT INTO product_variant (product_id, title, description, size, price, discounted, stock)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CreateVariantAttribute :one
INSERT INTO variant_attribute (variant_id, name, value)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateVariantImage :one
INSERT INTO variant_images (variant_id, image_key, position)
VALUES ($1, $2, $3)
RETURNING *;


-- name: UpdateProduct :one
UPDATE product
SET title = $2,
    description = $3,
    category_id = $4,
    is_active = $5
WHERE id = $1
RETURNING *;

-- name: UpdateProductVariant :one
UPDATE product_variant
SET title = $2,
    description = $3,
    size = $4,
    price = $5,
    discounted = $6,
    stock = $7
WHERE id = $1
RETURNING *;

-- name: UpdateVariantAttribute :one
UPDATE variant_attribute
SET name = $2,
    value = $3
WHERE id = $1
RETURNING *;  

-- name: UpdateVariantImage :one
UPDATE variant_images
SET image_key = $2,
    position = $3
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
UPDATE product
SET is_active = FALSE
WHERE id = $1;

-- name: GetProductVariant :one
SELECT * FROM product_variant
WHERE product_id = $1 and id = $2;