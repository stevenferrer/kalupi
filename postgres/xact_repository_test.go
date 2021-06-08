package postgres_test

import (
	"context"
	"testing"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/currency"
	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/ledger"
	"github.com/sf9v/kalupi/postgres"
	"github.com/sf9v/kalupi/transaction"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXactRepository(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	ctx := context.TODO()

	// setup accounts and ledgers
	accountRepo := postgres.NewAccountRepository(db)
	johnDoe := account.Account{
		AccountID: account.AccountID("johndoex"),
		Currency:  currency.USD,
	}
	_, err = accountRepo.CreateAccount(ctx, johnDoe)
	require.NoError(t, err)
	maryJane := account.Account{
		AccountID: account.AccountID("maryjanex"),
		Currency:  currency.USD,
	}
	_, err = accountRepo.CreateAccount(ctx, maryJane)
	require.NoError(t, err)

	ledgerRepo := postgres.NewLedgerRepository(db)
	ledgerService := ledger.NewService(ledgerRepo)
	err = ledgerService.CreateCashLedgers(ctx)
	require.NoError(t, err)

	cashUSDLedger := ledger.CashUSDLedgerNo

	balRepo := postgres.NewBalanceRepository(db)
	xactRepo := postgres.NewXactRepository(db)

	t.Run("create xact", func(t *testing.T) {
		t.Run("deposit", func(t *testing.T) {
			tx, err := xactRepo.BeginTx(ctx)
			require.NoError(t, err)

			// simulate deposit
			dpXactNo := transaction.XactNo("dp900")
			err = xactRepo.CreateXact(ctx, tx, transaction.Transaction{
				XactNo:      dpXactNo,
				LedgerNo:    cashUSDLedger,
				XactType:    transaction.XactTypeDebit,
				AccountID:   johnDoe.AccountID,
				XactTypeExt: transaction.XactTypeExtDeposit,
				Amount:      decimal.NewFromInt(100),
				Desc:        "Initial deposit",
			})
			require.NoError(t, err)

			err = tx.Commit()
			require.NoError(t, err)
		})

		t.Run("withdrawal", func(t *testing.T) {
			tx, err := xactRepo.BeginTx(ctx)
			require.NoError(t, err)

			// simulate withdrawal
			wdXactNo := transaction.XactNo("wd910")
			err = xactRepo.CreateXact(ctx, tx, transaction.Transaction{
				XactNo:      wdXactNo,
				LedgerNo:    cashUSDLedger,
				XactType:    transaction.XactTypeCredit,
				AccountID:   johnDoe.AccountID,
				XactTypeExt: transaction.XactTypeExtWithdrawal,
				Amount:      decimal.NewFromInt(25),
				Desc:        "Withdrawal",
			})
			require.NoError(t, err)

			err = tx.Commit()
			require.NoError(t, err)
		})

		t.Run("transfer", func(t *testing.T) {
			tx, err := xactRepo.BeginTx(ctx)
			require.NoError(t, err)

			trXactNo := transaction.XactNo("tr920")
			transferAmount := decimal.NewFromInt(25)
			err = xactRepo.CreateXact(ctx, tx, transaction.Transaction{
				XactNo:      trXactNo,
				LedgerNo:    cashUSDLedger,
				XactType:    transaction.XactTypeCredit,
				AccountID:   johnDoe.AccountID,
				XactTypeExt: transaction.XactTypeExtSndTransfer, // debit
				Amount:      transferAmount,
				Desc:        "Outgoing transfer to maryjane",
			})
			require.NoError(t, err)

			err = xactRepo.CreateXact(ctx, tx, transaction.Transaction{
				XactNo:      trXactNo,
				LedgerNo:    cashUSDLedger,
				XactType:    transaction.XactTypeDebit,
				AccountID:   maryJane.AccountID,
				XactTypeExt: transaction.XactTypeExtRcvTransfer, // credit
				Amount:      transferAmount,
				Desc:        "Incoming transfer from johndoe",
			})
			require.NoError(t, err)

			err = tx.Commit()
			require.NoError(t, err)
		})

		t.Run("verify balances", func(t *testing.T) {
			tx, err := xactRepo.BeginTx(ctx)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, tx.Commit())
			}()

			johnBal, err := balRepo.GetAccntBal(ctx, tx, johnDoe.AccountID)
			require.NoError(t, err)
			assert.True(t, decimal.NewFromInt(100).Equal(johnBal.TotalCredit))
			assert.True(t, decimal.NewFromInt(50).Equal(johnBal.TotalDebit))
			assert.True(t, decimal.NewFromInt(50).Equal(johnBal.CurrentBal))

			maryBal, err := balRepo.GetAccntBal(ctx, tx, maryJane.AccountID)
			require.NoError(t, err)
			assert.True(t, decimal.NewFromInt(25).Equal(maryBal.TotalCredit))
			assert.True(t, decimal.NewFromInt(0).Equal(maryBal.TotalDebit))
			assert.True(t, decimal.NewFromInt(25).Equal(maryBal.CurrentBal))
		})
	})

	t.Run("list xacts", func(t *testing.T) {
		xacts, err := xactRepo.ListXacts(ctx)
		require.NoError(t, err)
		assert.Len(t, xacts, 4)
	})
}
