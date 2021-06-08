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

type XactType int

const (
	XactTypeDebit XactType = iota + 1
	XactTypeCredit
)

func (tt XactType) String() string {
	return [...]string{
		"invalid",
		"Dr",
		"Cr",
	}[tt]
}

func (tt XactType) Value() (driver.Value, error) {
	return tt.String(), nil
}

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

func strToXactType(s string) XactType {
	switch s {
	case "Dr":
		return XactTypeDebit
	case "Cr":
		return XactTypeCredit
	}

	return XactType(0)
}

type XactTypeExt int

const (
	XactTypeExtDeposit XactTypeExt = iota + 1
	XactTypeExtWithdrawal
	// XactTypeExtSndTransfer is used when an account sends money to another account.
	// The sending account will be debited.
	XactTypeExtSndTransfer
	// XactTypeRcvTransfer is used when an account is receiving money from another account.
	// The recieving account will be credited.
	XactTypeExtRcvTransfer
)

func (ttx XactTypeExt) String() string {
	return [...]string{
		"invalid",
		"Dp",
		"Wd",
		"STr",
		"RTr",
	}[ttx]
}

func (ttx XactTypeExt) Value() (driver.Value, error) {
	return ttx.String(), nil
}

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
