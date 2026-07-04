-- name: CreateProdVariants :one
INSERT INTO product_variants (
  product_id,
  volume_ml,
  price,
  stock
) VALUES (
  $1, $2, $3, $4 
)RETURNING *;

-- name: GetProdVariantByProductID :many
SELECT * FROM product_variants 
WHERE product_id = $1 LIMIT 1;

-- name: GetProdVariantID :one
SELECT * FROM product_variants 
WHERE id = $1;

-- name: ListProdVariant :many
SELECT * FROM product_variants
ORDER BY product_id
LIMIT $1
OFFSET $2;

-- name: UpdateVariant :one
UPDATE product_variants
SET 
    volume_ml = COALESCE($2, volume_ml),
    price = COALESCE($3, price),
    stock = COALESCE($4, stock)
WHERE id = $1
RETURNING *;



-- name: UpdateProdVariantPrice :one
UPDATE product_variants
SET price = $2
WHERE id = $1
RETURNING *;

-- name: UpdateProdVariantStock :one
UPDATE product_variants
SET stock = $2
WHERE id = $1
RETURNING *;

-- name: UpdateProdVariantVolume :one
UPDATE product_variants
SET volume_ml = $2
WHERE id = $1
RETURNING *;

-- name: DecreaseProdVariantStock :one
UPDATE product_variants 
SET stock = stock - $2
WHERE id = $1 AND stock >= $2
RETURNING *;


-- name: DeleteProdVariant :exec
DELETE FROM product_variants
WHERE id = $1;
