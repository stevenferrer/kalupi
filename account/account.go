package account

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/shopspring/decimal"

	"github.com/sf9v/kalupi/currency"
)

// AccountID is the account id
type AccountID string

// Validate validates the account id
func (accntID AccountID) Validate() error {
	return validation.Validate(string(accntID),
		validation.Required.Error("must not be empty"),
		validation.Length(6, 64).Error("must have length between 6 and 64"),
		is.Alphanumeric.Error("must contain english letters and digits only"),
	)
}

// Account is an external account. These accounts are usually debit accounts.
type Account struct {
	AccountID AccountID         `json:"id"`
	Currency  currency.Currency `json:"currency"`
	Balance   decimal.Decimal   `json:"balance"`
}

// Validate validates the account
func (ac Account) Validate() error {
	return validation.Errors{
		"account_id": ac.AccountID.Validate(),
		"currency": validation.Validate(ac.Currency,
			validation.Required.Error("currency is required"),
			validation.By(func(value interface{}) error {
				c, _ := value.(currency.Currency)
				if !c.IsValid() {
					return currency.ErrUnsupportedCurrency
				}

				return nil
			})),
	}.Filter()
}

// Balance is an account balance
type Balance struct {
	AccountID   AccountID
	TotalCredit decimal.Decimal
	TotalDebit  decimal.Decimal
	CurrentBal  decimal.Decimal
	Ts          *time.Time
}
