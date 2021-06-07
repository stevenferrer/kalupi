package account_test

import (
	"testing"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/currency"
	"github.com/stretchr/testify/assert"
)

func TestAccountValidation(t *testing.T) {
	t.Run("empty account id", func(t *testing.T) {
		ac := account.Account{
			AccountID: "",
			Currency:  currency.USD,
		}
		assert.Error(t, ac.Validate())
	})

	t.Run("account id max letters", func(t *testing.T) {
		ac := account.Account{
			AccountID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa---",
			Currency:  currency.USD,
		}
		assert.Error(t, ac.Validate())
	})

	t.Run("currency not supported", func(t *testing.T) {
		ac := account.Account{AccountID: "john123"}
		assert.Error(t, ac.Validate())
	})
}
