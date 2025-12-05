-- name: CreateUser :one
INSERT INTO "user" (email, role, username)
VALUES ($1, $2, $3) 
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

-- Product Queries

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


-- name: GetUploadRequestByKey :one
SELECT * FROM uploads
WHERE key = $1;

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

-- name: GetProductBySeller :one
SELECT * FROM product
WHERE id = $1 AND seller_id = $2;

-- name: GetSellerByUserID :one
SELECT * FROM seller
WHERE user_id = $1;

-- name: CreateRequestFileUpload :one
INSERT INTO uploads (key, filename, content_type, file_size, upload_type, expires_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;