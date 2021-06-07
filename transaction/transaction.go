package transaction

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/ledger"
)

// XactNo is a transaction number
type XactNo string

// Transaction represents a double-entry record for an external account transaction.
// Note: Ledger and account must have the same currency.
type Transaction struct {
	// XactNo is a transaction reference number
	XactNo XactNo

	LedgerNo ledger.LedgerNo
	XactType XactType

	AccountID   account.AccountID
	XactTypeExt XactTypeExt

	// Amount is the amount of transaction
	Amount decimal.Decimal

	// Desc is a short description of entry i.e. deposit, withdrawal
	Desc string

	Ts *time.Time
}
