-- name: CreateUser :exec
INSERT INTO users (
    id, username, password
) VALUES (?, ?, ?);

-- name: UserExists :one
SELECT EXISTS (
    SELECT id FROM users WHERE username = ?
);
