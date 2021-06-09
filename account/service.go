package account

import (
	"context"
)

// Service is an account service
type Service interface {
	// CreateAccount creates an account
	CreateAccount(context.Context, Account) error
	// GetAccount retrieives an account via AccountID
	GetAccount(context.Context, AccountID) (*Account, error)
	// ListAccounts retrieives the list of accounts
	ListAccounts(context.Context) ([]*Account, error)
}
