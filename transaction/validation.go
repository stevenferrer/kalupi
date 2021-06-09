package transaction

import (
	"github.com/shopspring/decimal"
)

// nonZeroDecimal validates the decimal is non-zero
func nonZeroDecimal(value interface{}) error {
	c, _ := value.(decimal.Decimal)
	if c.IsZero() {
		return ErrZeroAmount
	}

	return nil
}

// nonNegativeDecimal validates the decimal is non-negative
func nonNegativeDecimal(value interface{}) error {
	c, _ := value.(decimal.Decimal)
	if c.IsNegative() {
		return ErrNegativeAmount
	}

	return nil
}
