package transaction

import (
	"context"

	"github.com/sf9v/kalupi/etc/tx"
)

// Repository is a transaction repository
type Repository interface {
	// BeginTx begins a new tx
	BeginTx(context.Context) (tx.Tx, error)
	// CreateXact creates a transaction within tx
	CreateXact(context.Context, tx.Tx, Transaction) error
	// ListXacts retrieves the list of transactions
	ListXacts(context.Context) ([]*Transaction, error)
	// ListTransfers retrieves the transfer related transactions
	ListTransfers(context.Context) ([]*Transaction, error)
}
