package ledger

import (
	"database/sql/driver"
	"errors"
)

// AccountType is a ledger account type
type AccountType int

// List of ledger account types
const (
	// AccountTypeLiability is a liability account
	AccountTypeLiability AccountType = iota + 1
)

// String implements Stringer interface
func (at AccountType) String() string {
	return [...]string{
		"invalid",
		"AL",
	}[at]
}

// Value implements the driver.Valuer interface
func (at AccountType) Value() (driver.Value, error) {
	return at.String(), nil
}

// Scan implements the sql.Scanner interface
func (at *AccountType) Scan(src interface{}) error {
	if src == nil {
		*at = AccountType(0)
		return nil
	}

	val, ok := src.(string)
	if !ok {
		return errors.New("src is not string")
	}

	*at = strToAccountType(val)
	return nil
}

// strToAccountType takes a string and returns the account type
func strToAccountType(s string) AccountType {
	switch s {
	case "AL":
		return AccountTypeLiability
	}

	return AccountType(0)
}
