package transaction

import (
	"database/sql/driver"
	"errors"
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

	// LedgerNo is the ledger number of internal account
	LedgerNo ledger.LedgerNo
	// XactType is the transaction type (debit or credit)
	XactType XactType

	// AccountId is the external account id
	AccountID account.AccountID
	// XactTypeExt is the external transaction type (deposit, withdrawal, transfer)
	XactTypeExt XactTypeExt

	// Amount is the amount of transaction
	Amount decimal.Decimal

	// Desc is a short description of entry i.e. deposit, withdrawal
	Desc string

	// Ts is the timestamp
	Ts *time.Time
}

// XactType is the transaction type
type XactType int

// List of transaction types
const (
	// XactTypeDebit is a debit transaction
	XactTypeDebit XactType = iota + 1
	// XactTypeCredit is a credit transaction
	XactTypeCredit
)

// String implements Stringer
func (tt XactType) String() string {
	return [...]string{
		"invalid",
		"Dr",
		"Cr",
	}[tt]
}

// Value implements driver.Valuer interface
func (tt XactType) Value() (driver.Value, error) {
	return tt.String(), nil
}

// Scan implements sql.Scanner interface
func (tt *XactType) Scan(src interface{}) error {
	if src == nil {
		*tt = XactType(0)
		return nil
	}

	val, ok := src.(string)
	if !ok {
		return errors.New("src is not string")
	}

	*tt = strToXactType(val)
	return nil
}

// strToXactType takes a string and returns the transaction type
func strToXactType(s string) XactType {
	switch s {
	case "Dr":
		return XactTypeDebit
	case "Cr":
		return XactTypeCredit
	}

	return XactType(0)
}

// XactTypeExt is an external transaction type
type XactTypeExt int

// List of external transaction types
const (
	// XactTypeExt is a deposit transaction
	XactTypeExtDeposit XactTypeExt = iota + 1
	// XactTypeExtWithdrawal is a withdrawal transaction
	XactTypeExtWithdrawal
	// XactTypeExtSndTransfer is a outgoing transfer.
	// The sending account will be debited.
	XactTypeExtSndTransfer
	// XactTypeRcvTransfer is an incomming transfer.
	// The receving account will be credited.
	XactTypeExtRcvTransfer
)

// String implements Stringer interface
func (ttx XactTypeExt) String() string {
	return [...]string{
		"invalid",
		"Dp",
		"Wd",
		"STr",
		"RTr",
	}[ttx]
}

// Value implements driver.Valuer interface
func (ttx XactTypeExt) Value() (driver.Value, error) {
	return ttx.String(), nil
}

// Scan implements sql.Scanner interface
func (ttx *XactTypeExt) Scan(src interface{}) error {
	if src == nil {
		*ttx = XactTypeExt(0)
		return nil
	}

	val, ok := src.(string)
	if !ok {
		return errors.New("src is not string")
	}

	*ttx = strToXactTypeExt(val)
	return nil
}

// strToXactTypeExt takes a string and returns the external transaction type
func strToXactTypeExt(s string) XactTypeExt {
	switch s {
	case "Dp":
		return XactTypeExtDeposit
	case "Wd":
		return XactTypeExtWithdrawal
	case "STr":
		return XactTypeExtSndTransfer
	case "RTr":
		return XactTypeExtRcvTransfer
	}

	return XactTypeExt(0)
}
