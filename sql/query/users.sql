-- name: CreateUser :exec
INSERT INTO users (
    id, username, password
) VALUES (?, ?, ?);
