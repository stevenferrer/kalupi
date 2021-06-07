package ledger

import (
	"database/sql/driver"
	"errors"
)

type AccountType int

const (
	AccountTypeLiability AccountType = iota + 1
)

func (at AccountType) String() string {
	return [...]string{
		"invalid",
		"AL",
	}[at]
}

func (at AccountType) Value() (driver.Value, error) {
	return at.String(), nil
}

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

func strToAccountType(s string) AccountType {
	switch s {
	case "AL":
		return AccountTypeLiability
	}

	return AccountType(0)
}
