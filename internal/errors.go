package internal

import (
	"net/http"
)

// MakeBadRequestError

type BadRequestError struct{ error error }

func (e BadRequestError) Error() string {
	return e.error.Error()
}

func (e BadRequestError) Code() int {
	return http.StatusBadRequest
}

func MakeBadRequestError(error error) BadRequestError {
	return BadRequestError{error: error}
}

// Validation

type ValidationError struct{ error error }

func (e ValidationError) Error() string {
	return e.error.Error()
}

func (e ValidationError) Code() int {
	return http.StatusUnprocessableEntity
}

func MakeValidationError(error error) ValidationError {
	return ValidationError{error: error}
}

// Authentication

type AuthenticationError struct{ error error }

func (e AuthenticationError) Error() string {
	return e.error.Error()
}

func (e AuthenticationError) Code() int {
	return http.StatusUnauthorized
}

func MakeAuthenticationError(error error) AuthenticationError {
	return AuthenticationError{error: error}
}
