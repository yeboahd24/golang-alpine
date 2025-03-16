-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CheckEmailExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE email = $1
) AS exists;

-- name: CheckUsernameExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE username = $1
) AS exists;