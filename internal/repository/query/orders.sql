-- name: CreateOrder :one
INSERT INTO orders (
  user_id,
  product_id,
  quantity,
  total
) VALUES (
  $1, $2, $3, $4
)RETURNING *;


-- name: GetOrderByUser :one
SELECT * FROM orders
WHERE user_id = $1 LIMIT 1;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrderByID :many
SELECT * FROM orders
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: ListOrderByTime :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;


  -- name: UpdateOrderTotal :one
  UPDATE orders
  SET total = $2
  WHERE id = $1
  RETURNING *;


-- name: DeleteOrderByID :exec
DELETE FROM orders
WHERE id = $1;


