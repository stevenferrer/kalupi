package account

import "errors"

// List of account related errors
var (
	ErrAccountAlreadyExist = errors.New("account already exist")
	ErrAccountNotFound     = errors.New("account not found")
)
