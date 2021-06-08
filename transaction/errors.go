package transaction

import "errors"

// List of transaction related errors
var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrDifferentCurrencies = errors.New("different currencies")
)
