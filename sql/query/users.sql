-- name: CreateUser :exec
INSERT INTO users (
    id, username, password
) VALUES ($1, $2, $3);

-- name: UserExists :one
SELECT EXISTS (
    SELECT id FROM users WHERE username = $1
);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users where id = $1;

-- name: UpdateUserEmail :exec
UPDATE users SET email = $1 WHERE id = $2;
