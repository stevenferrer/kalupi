package postgres_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/currency"
	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/ledger"
	"github.com/sf9v/kalupi/postgres"
	"github.com/sf9v/kalupi/transaction"
)

func TestBalanceRepository(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	ctx := context.TODO()

	// setup accounts and ledgers
	accountRepo := postgres.NewAccountRepository(db)
	john := account.Account{
		AccountID: account.AccountID("johndoe"),
		Currency:  currency.USD,
	}
	_, err = accountRepo.CreateAccount(ctx, john)
	require.NoError(t, err)
	mary := account.Account{
		AccountID: account.AccountID("maryjane"),
		Currency:  currency.USD,
	}
	_, err = accountRepo.CreateAccount(ctx, mary)
	require.NoError(t, err)

	ledgerRepo := postgres.NewLedgerRepository(db)
	cashUSD := ledger.Ledger{
		LedgerNo:    ledger.CashUSDLedgerNo,
		AccountType: ledger.AccountTypeLiability,
		Currency:    currency.USD,
		Name:        "Cash USD",
	}
	err = ledgerRepo.CreateLedgersIfNotExist(ctx, cashUSD)
	require.NoError(t, err)

	balRepo := postgres.NewBalanceRepository(db)
	xactRepo := postgres.NewXactRepository(db)

	// seed some money
	{
		tx1, err := xactRepo.BeginTx(ctx)
		require.NoError(t, err)

		// john should have 0 balance
		johnBal, err := balRepo.GetAccntBal(ctx, tx1, john.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(0).Equal(johnBal.CurrentBal),
			"john should have an initial balance of 0")

		dpXactNo, err := transaction.NewXactNo()
		require.NoError(t, err)
		err = xactRepo.CreateXact(ctx, tx1, transaction.Transaction{
			XactNo:      dpXactNo,
			LedgerNo:    cashUSD.LedgerNo,
			XactType:    transaction.XactTypeDebit,
			AccountID:   john.AccountID,
			XactTypeExt: transaction.XactTypeExtDeposit,
			Amount:      decimal.NewFromInt(100),
			Desc:        "Cash deposit from johndoe",
		})
		require.NoError(t, err)
		require.NoError(t, tx1.Commit())
	}

	t.Run("concurrent access", func(t *testing.T) {
		chanGetBal := make(chan bool)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			tx, err := xactRepo.BeginTx(ctx)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, tx.Commit())
			}()

			// simulate withdrawal
			wdXactNo, err := transaction.NewXactNo()
			require.NoError(t, err)
			err = xactRepo.CreateXact(ctx, tx, transaction.Transaction{
				XactNo:      wdXactNo,
				LedgerNo:    cashUSD.LedgerNo,
				XactType:    transaction.XactTypeCredit,
				AccountID:   john.AccountID,
				XactTypeExt: transaction.XactTypeExtWithdrawal,
				Amount:      decimal.NewFromInt(100),
				Desc:        "Cash withdrawal from johndoe",
			})
			require.NoError(t, err)

			// this will ensure create xact comes first before the balance check
			chanGetBal <- true
			close(chanGetBal)

			// block and let others try to read balance
			time.Sleep(10 * time.Millisecond)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			<-chanGetBal

			tx, err := balRepo.BeginTx(ctx)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, tx.Commit())
			}()

			// ideally, this should block and wait for the above code!
			johnBal, err := balRepo.GetAccntBal(ctx, tx, john.AccountID)
			require.NoError(t, err)

			assert.True(t, decimal.NewFromInt(0).Equal(johnBal.CurrentBal),
				"john should have 0 balance after withdrawal")
		}()

		wg.Wait()
	})
}
