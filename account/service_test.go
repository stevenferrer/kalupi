package account_test

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

func TestAccountService(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	accntRepo := postgres.NewAccountRepository(db)
	accntService := account.NewService(accntRepo)

	ctx := context.TODO()
	accountID := account.AccountID("john1234")
	t.Run("create account", func(t *testing.T) {
		err := accntService.CreateAccount(ctx, account.Account{
			AccountID: accountID,
			Currency:  currency.USD,
		})
		require.NoError(t, err)
	})

	t.Run("get account", func(t *testing.T) {
		ac, err := accntService.GetAccount(ctx, accountID)
		require.NoError(t, err)

		assert.Equal(t, accountID, ac.AccountID)
		assert.Equal(t, currency.USD, ac.Currency)
	})

	t.Run("list accounts", func(t *testing.T) {
		acs, err := accntService.ListAccounts(ctx)
		require.NoError(t, err)

		assert.Len(t, acs, 1)
	})

}
