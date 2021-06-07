package ledger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/ledger"
	"github.com/sf9v/kalupi/postgres"
)

func TestLedgerService(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	ledgerRepo := postgres.NewLedgerRepository(db)
	ledgerService := ledger.NewService(ledgerRepo)

	ctx := context.TODO()
	t.Run("create cash ledgers", func(t *testing.T) {
		err := ledgerService.CreateCashLedgers(ctx)
		require.NoError(t, err)

		lgs, err := ledgerRepo.ListLedgers(ctx)
		require.NoError(t, err)
		assert.Len(t, lgs, 1)
	})
}
