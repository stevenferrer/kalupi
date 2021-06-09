package transaction

import "errors"

// List of transaction related errors
var (
	// ErrInsufficientBalance is an error when  the
	// account balance is insufficient for the transaction
	ErrInsufficientBalance = errors.New("insufficient balance")
	// ErrDifferentCurrencies is an error when transferring money
	// to an account with different currency
	ErrDifferentCurrencies = errors.New("different currencies")
	// ErrSendingAccountNotFound is an error when
	// transfering money from an account that doesn't exist
	ErrSendingAccountNotFound = errors.New("sending account not found")
	// ErrReceivingAccountNotFound is an error when
	// 	transfering money to an account that doesn't exist
	ErrReceivingAccountNotFound = errors.New("receiving account not found")
	// ErrValidation is a transaction related validation error
	ErrValidation = errors.New("validation error")
	// ErrZeroAmount is an error when doing a transaction with a zero amount
	ErrZeroAmount = errors.New("zero amount")
	// ErrNegativeAmount is an error when doing a transaction with negative amount
	ErrNegativeAmount = errors.New("negative amount")
)
