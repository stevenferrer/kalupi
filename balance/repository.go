package balance

import (
	"context"

	"github.com/stevenferrer/kalupi/account"
	"github.com/stevenferrer/kalupi/etc/tx"
)

// Repository is a balance repository
type Repository interface {
	// BeginTx begins a new transaction
	BeginTx(context.Context) (tx.Tx, error)
	// GetAccntBal retreives the account balance within tx
	GetAccntBal(context.Context, tx.Tx, account.AccountID) (*account.Balance, error)
}
