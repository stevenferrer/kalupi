package account

import (
	"errors"
	"time"

	"go.uber.org/multierr"

	"github.com/sf9v/kalupi/currency"
	"github.com/shopspring/decimal"
)

// TODO: validations for new accounts

// AccountID is the account id
type AccountID string

func (acID AccountID) Validate() error {
	var errs error
	if len(acID) == 0 {
		errs = multierr.Append(errs, errors.New("account id cannot be empty"))
	}

	if len(acID) > 64 {
		errs = multierr.Append(errs, errors.New("account id can only have 64 letters"))
	}

	return errs
}

// Account is an external customer account.
// These are usually debit accounts.
type Account struct {
	AccountID AccountID
	Currency  currency.Currency
}

func (ac Account) Validate() error {
	var errs error

	if err := ac.AccountID.Validate(); err != nil {
		errs = multierr.Append(errs, err)
	}

	if !ac.Currency.IsValid() {
		errs = multierr.Append(errs, errors.New("currency not supported"))
	}

	return errs
}

// Balance is the account balance
type Balance struct {
	AccntID     AccountID
	TotalCredit decimal.Decimal
	TotalDebit  decimal.Decimal
	CurrentBal  decimal.Decimal
	Ts          *time.Time
}
