package repositories

import (
	"context"
	"database/sql"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal/db"
)

func UserExists(ctx context.Context, username string) (bool, error) {
	return db.New(database.DB).UserExists(ctx, username)
}

func CreateUser(ctx context.Context, id string, username string, hashedPassword string) error {
	return db.New(database.DB).CreateUser(ctx, db.CreateUserParams{
		ID:       id,
		Username: username,
		Password: hashedPassword,
	})
}

func GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return db.New(database.DB).GetUserByUsername(ctx, username)
}

func GetUserByID(ctx context.Context, userID string) (db.User, error) {
	return db.New(database.DB).GetUserByID(ctx, userID)
}

func UpdateUserEmail(ctx context.Context, userID string, email string) error {
	return db.New(database.DB).UpdateUserEmail(ctx, db.UpdateUserEmailParams{
		Email: sql.NullString{String: email, Valid: true},
		ID:    userID,
	})
}
