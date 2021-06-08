package account

import (
	"context"
)

// Repository is an account repository
type Repository interface {
	// CreateAccount creates an account record
	CreateAccount(context.Context, Account) (AccountID, error)
	// GetAccount retrieves the account record using account id
	GetAccount(context.Context, AccountID) (*Account, error)
	// IsAccountExists returns true if the account id exists
	IsAccountExists(context.Context, AccountID) (bool, error)
	// ListAccounts retrieves the account records
	ListAccounts(context.Context) ([]*Account, error)
}
