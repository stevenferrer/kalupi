package account

import (
	"context"
)

// Repository is an account repository
type Repository interface {
	// CreateAccount creates an account
	CreateAccount(context.Context, Account) (AccountID, error)
	// GetAccount retrieves the account
	GetAccount(context.Context, AccountID) (*Account, error)
	// IsAccountExists returns true if the account exists
	IsAccountExists(context.Context, AccountID) (bool, error)
	// ListAccounts retrieves the list of accounts
	ListAccounts(context.Context) ([]*Account, error)
}
