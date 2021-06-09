package currency

import "errors"

// List of currency related errors
var (
	// ErrUnsupportedCurrency is an error when a currency is not supported
	ErrUnsupportedCurrency = errors.New("unsupported currency")
)
