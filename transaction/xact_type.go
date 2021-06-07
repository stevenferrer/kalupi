package transaction

import (
	"database/sql/driver"
	"errors"
)

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
