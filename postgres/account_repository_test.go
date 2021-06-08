package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/currency"
	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/postgres"
)

func TestAccountRepository(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	accountRepo := postgres.NewAccountRepository(db)

	accountID := account.AccountID("john1234")
	ctx := context.TODO()
	t.Run("create account", func(t *testing.T) {
		_, err = accountRepo.CreateAccount(ctx, account.Account{
			AccountID: accountID,
			Currency:  currency.USD,
		})
		require.NoError(t, err)
	})

	t.Run("get account", func(t *testing.T) {
		a, err := accountRepo.GetAccount(ctx, accountID)
		require.NoError(t, err)
		assert.Equal(t, accountID, a.AccountID)
		assert.Equal(t, currency.USD, a.Currency)
	})

	t.Run("list accounts", func(t *testing.T) {
		as, err := accountRepo.ListAccounts(ctx)
		require.NoError(t, err)
		assert.Len(t, as, 1)
	})

	t.Run("account exists", func(t *testing.T) {
		exists, err := accountRepo.IsAccountExists(ctx, accountID)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = accountRepo.IsAccountExists(ctx, account.AccountID("idontexist"))
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
