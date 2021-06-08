package transaction

import (
	"context"

	"github.com/sf9v/kalupi/etc/tx"
)

type Repository interface {
	BeginTx(context.Context) (tx.Tx, error)
	CreateXact(context.Context, tx.Tx, Transaction) error
	ListXacts(context.Context) ([]*Transaction, error)
	ListTransfers(context.Context) ([]*Transaction, error)
}
