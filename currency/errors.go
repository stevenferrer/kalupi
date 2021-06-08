package currency

import "errors"

// List of currency related errors
var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
)
