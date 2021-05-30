package services

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/ledger-auth/internal"
	"github.com/bradenrayhorn/ledger-auth/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username string, password string) error {
	exists, err := repositories.UserExists(context.Background(), username)
	if err != nil {
		return internal.MakeBadRequestError(err)
	}

	if exists {
		return internal.MakeValidationError(errors.New("user already exists"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return internal.MakeBadRequestError(err)
	}

	err = repositories.CreateUser(context.Background(), uuid.NewString(), username, string(hashedPassword))
	if err != nil {
		return internal.MakeBadRequestError(err)
	}

	return nil
}

func Login(username string, password string) (string, error) {
	user, err := repositories.GetUserByUsername(context.Background(), username)
	if err != nil {
		return "", internal.MakeAuthenticationError(errors.New("invalid username/password"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", internal.MakeAuthenticationError(errors.New("invalid username/password"))
	}

	return user.ID, nil
}
