package account_test

import (
	"testing"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/currency"
	"github.com/stretchr/testify/assert"
)

func TestAccountValidation(t *testing.T) {
	t.Run("account id", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			ac := account.Account{AccountID: "", Currency: currency.USD}
			err := ac.Validate()
			assert.Error(t, err)
		})

		t.Run("min and max letters", func(t *testing.T) {
			ac := account.Account{AccountID: "john", Currency: currency.USD}
			err := ac.Validate()
			assert.Error(t, err)

			ac = account.Account{
				AccountID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa---",
				Currency:  currency.USD,
			}
			err = ac.Validate()
			assert.Error(t, err)
		})

		t.Run("alphanumeric letters only", func(t *testing.T) {
			ac := account.Account{AccountID: "john#$#$#$", Currency: currency.USD}
			err := ac.Validate()
			assert.Error(t, err)
		})
	})

	t.Run("currency not supported", func(t *testing.T) {
		ac := account.Account{AccountID: "johnybravo"}
		err := ac.Validate()
		assert.Error(t, err)
	})

}
