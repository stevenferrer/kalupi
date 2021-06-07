package transaction

import (
	"database/sql/driver"
	"errors"
)

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
