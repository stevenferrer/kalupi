package transaction

import (
	"github.com/shopspring/decimal"
)

func nonZeroDecimal(value interface{}) error {
	c, _ := value.(decimal.Decimal)
	if c.IsZero() {
		return ErrZeroAmount
	}

	return nil
}

func nonNegativeDecimal(value interface{}) error {
	c, _ := value.(decimal.Decimal)
	if c.IsNegative() {
		return ErrNegativeAmount
	}

	return nil
}
