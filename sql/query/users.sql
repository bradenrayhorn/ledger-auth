-- name: CreateUser :exec
INSERT INTO users (
    id, username, password
) VALUES (?, ?, ?);

-- name: UserExists :one
SELECT EXISTS (
    SELECT id FROM users WHERE username = ?
);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUserByID :one
SELECT * FROM users where id = ?;

-- name: UpdateUserEmail :exec
UPDATE users SET email = ? WHERE id = ?;
