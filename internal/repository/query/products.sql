-- name: CreateProduct :one
INSERT INTO products (
  name,
  price,
  stock,
  description
) VALUES (
  $1, $2, $3, $4
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
ORDER BY name
LIMIT $1
OFFSET $2;


-- name: UpdateProductPrice :one
UPDATE products
SET price = $2
WHERE id = $1
RETURNING *;

-- name: UpdateProductStock :one
UPDATE products
SET stock = $2
WHERE id = $1
RETURNING *;
