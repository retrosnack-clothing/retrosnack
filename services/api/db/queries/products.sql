-- name: ListProducts :many
SELECT id, title, description, category_id, brand, condition,
       price_cents, seller_id, instagram_post_url, created_at, updated_at
FROM products
ORDER BY created_at DESC;

-- name: GetProduct :one
SELECT id, title, description, category_id, brand, condition,
       price_cents, seller_id, instagram_post_url, created_at, updated_at
FROM products
WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (title, description, category_id, brand, condition, price_cents, instagram_post_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
