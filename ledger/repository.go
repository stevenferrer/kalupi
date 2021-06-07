package ledger

import "context"

type Repository interface {
	CreateLedgersIfNotExist(context.Context, ...Ledger) error
	GetLedger(context.Context, LedgerNo) (*Ledger, error)
	ListLedgers(context.Context) ([]*Ledger, error)
}
