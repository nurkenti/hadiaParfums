-- name: CreateUser :one
INSERT INTO users (
  chat_id,
  lang
  
) VALUES (
  $1, $2 
)RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE chat_id = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM users
ORDER BY chat_id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
SET lang = $2
WHERE chat_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE chat_id = $1;
