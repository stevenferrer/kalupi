package balance

import (
	"context"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/etc/tx"
)

type Repository interface {
	BeginTx(context.Context) (tx.Tx, error)
	GetAccntBal(context.Context, tx.Tx, account.AccountID) (*account.Balance, error)
}
