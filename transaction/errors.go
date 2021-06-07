package transaction

import "errors"

var (
	ErrInsufficientBalance  = errors.New("insiffucient balance")
	ErrMustHaveSameCurrency = errors.New("must have same currency")
)
