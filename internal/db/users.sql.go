// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
    id, username, password
) VALUES (?, ?, ?)
`

type CreateUserParams struct {
	ID       string
	Username string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.ID, arg.Username, arg.Password)
	return err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, password, created_at, updated_at FROM users WHERE username = ?
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userExists = `-- name: UserExists :one
SELECT EXISTS (
    SELECT id FROM users WHERE username = ?
)
`

func (q *Queries) UserExists(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRowContext(ctx, userExists, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
