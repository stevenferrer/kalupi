package postgres_test

import (
	"context"
	"testing"

	"github.com/stevenferrer/kalupi/currency"
	"github.com/stevenferrer/kalupi/etc/txdb"
	"github.com/stevenferrer/kalupi/ledger"
	"github.com/stevenferrer/kalupi/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerRepository(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	ledgerRepo := postgres.NewLedgerRepository(db)

	ctx := context.Background()
	t.Run("create ledger", func(t *testing.T) {
		cashLedger := ledger.Ledger{
			LedgerNo:    ledger.CashUSDLedgerNo,
			AccountType: ledger.AccountTypeLiability,
			Currency:    currency.USD,
			Name:        "Cash USD",
		}
		// create ledger
		err := ledgerRepo.CreateLedgersIfNotExists(ctx, cashLedger)
		require.NoError(t, err)
		// create ledger (again), this should do nothing since we already
		// created the ledger above
		err = ledgerRepo.CreateLedgersIfNotExists(ctx, cashLedger)
		require.NoError(t, err)
	})

	t.Run("get ledger", func(t *testing.T) {
		lg, err := ledgerRepo.GetLedger(ctx, ledger.CashUSDLedgerNo)
		require.NoError(t, err)

		assert.Equal(t, ledger.CashUSDLedgerNo, lg.LedgerNo)
		assert.Equal(t, ledger.AccountTypeLiability, lg.AccountType)
		assert.Equal(t, currency.USD, lg.Currency)
		assert.Equal(t, "Cash USD", lg.Name)
	})

	t.Run("list ledgers", func(t *testing.T) {
		lgs, err := ledgerRepo.ListLedgers(ctx)
		require.NoError(t, err)
		assert.Len(t, lgs, 1)
	})
}
