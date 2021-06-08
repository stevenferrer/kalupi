package transaction

import "errors"

// List of transaction related errors
var (
	ErrInsufficientBalance      = errors.New("insufficient balance")
	ErrDifferentCurrencies      = errors.New("different currencies")
	ErrSendingAccountNotFound   = errors.New("sending account not found")
	ErrReceivingAccountNotFound = errors.New("receiving account not found")
	ErrValidation               = errors.New("validation error")
	ErrZeroAmount               = errors.New("zero amount")
	ErrNegativeAmount           = errors.New("negative amount")
)
