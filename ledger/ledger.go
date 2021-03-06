package ledger

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/stevenferrer/kalupi/currency"
)

// LedgerNo is a ledger account number
type LedgerNo string

// Ledger is an internal ledger account
type Ledger struct {
	LedgerNo    LedgerNo
	AccountType AccountType       // i.e. Liability
	Currency    currency.Currency // i.e. USD
	Name        string            // i.e. Cash - USD
}

// Balance is a ledger balance
type Balance struct {
	LedgerNo    LedgerNo
	TotalCredit decimal.Decimal
	TotalDebit  decimal.Decimal
	CurrentBal  decimal.Decimal
	Ts          *time.Time
}
