-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;

-- name: GetAllUsers :many
SELECT * FROM users;