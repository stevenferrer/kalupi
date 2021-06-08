package account

import (
	"context"
	"errors"
)

var (
	ErrValidation = errors.New("validation error")
)

// Service is an account service
type Service interface {
	CreateAccount(context.Context, Account) error
	GetAccount(context.Context, AccountID) (*Account, error)
	ListAccounts(context.Context) ([]*Account, error)
}
