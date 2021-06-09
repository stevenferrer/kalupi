package account

import "errors"

// List of account related errors
var (
	// ErrAccountAlreadyExists is an error when creating an account that already exists
	ErrAccountAlreadyExists = errors.New("account already exists")
	// ErrAccountNotFound is an error when retrieving an account that doesn't exists
	ErrAccountNotFound = errors.New("account not found")
	// ErrValidation is an account related validation error
	ErrValidation = errors.New("validation error")
)
