package ledger

import (
	"context"
)

// Repository is a ledger repository
type Repository interface {
	// CreateLedgersIfNotExists
	CreateLedgersIfNotExists(context.Context, ...Ledger) error
	// GetLedger retrieves the ledger
	GetLedger(context.Context, LedgerNo) (*Ledger, error)
	// ListLedgers retrieves the list of ledgers
	ListLedgers(context.Context) ([]*Ledger, error)
}
