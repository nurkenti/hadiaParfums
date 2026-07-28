-- name: CreateProduct :one
INSERT INTO products (
  name,
  category,
  description
) VALUES (
  $1, $2, $3 
)RETURNING *;


-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProductByName :one
SELECT * FROM products
WHERE name = $1 LIMIT 1;

-- name: ListProductByID :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListProductByName :many
SELECT * FROM products
WHERE name ILIKE $1
ORDER BY name
LIMIT $2
OFFSET $3;


-- name: UpdateProductCategory :one
UPDATE products
SET category = $2
WHERE id = $1
RETURNING *;

-- name: UpdateProductDiscription :one
UPDATE products
SET description = $2
WHERE id = $1
RETURNING *;

-- name: DeleteProductByID :exec
DELETE FROM products
WHERE id = $1;
