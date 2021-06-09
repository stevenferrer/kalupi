package service_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/account"
	accountservice "github.com/sf9v/kalupi/account/service"
	"github.com/sf9v/kalupi/balance"
	"github.com/sf9v/kalupi/currency"
	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/postgres"
)

func TestAccountService(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	balRepo := postgres.NewBalanceRepository(db)
	balService := balance.NewService(balRepo)

	accntRepo := postgres.NewAccountRepository(db)
	accountSvc := accountservice.New(accntRepo, balService)

	ctx := context.TODO()
	accountID := account.AccountID("john1234")
	t.Run("create account", func(t *testing.T) {
		err := accountSvc.CreateAccount(ctx, account.Account{
			AccountID: accountID,
			Currency:  currency.USD,
		})
		require.NoError(t, err)

		t.Run("validation error", func(t *testing.T) {
			err = accountSvc.CreateAccount(ctx, account.Account{})
			assert.ErrorIs(t, err, account.ErrValidation)
		})

		t.Run("alread exists", func(t *testing.T) {
			err = accountSvc.CreateAccount(ctx, account.Account{
				AccountID: accountID,
				Currency:  currency.USD,
			})
			assert.ErrorIs(t, err, account.ErrAccountAlreadyExists)
		})
	})

	t.Run("get account", func(t *testing.T) {
		ac, err := accountSvc.GetAccount(ctx, accountID)
		require.NoError(t, err)

		assert.Equal(t, accountID, ac.AccountID)
		assert.Equal(t, currency.USD, ac.Currency)
		assert.True(t, decimal.Zero.Equal(ac.Balance))

		t.Run("validation error", func(t *testing.T) {
			ac, err = accountSvc.GetAccount(ctx, account.AccountID("##(#*"))
			assert.ErrorIs(t, err, account.ErrValidation)
			assert.Nil(t, ac)
		})

		t.Run("not found", func(t *testing.T) {
			ac, err = accountSvc.GetAccount(ctx, account.AccountID("johntravolta"))
			assert.ErrorIs(t, err, account.ErrAccountNotFound)
			assert.Nil(t, ac)
		})
	})

	t.Run("list accounts", func(t *testing.T) {
		acs, err := accountSvc.ListAccounts(ctx)
		require.NoError(t, err)

		assert.Len(t, acs, 1)
	})

}
