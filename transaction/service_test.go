package transaction_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/balance"
	"github.com/sf9v/kalupi/currency"
	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/ledger"
	"github.com/sf9v/kalupi/postgres"
	"github.com/sf9v/kalupi/transaction"
)

func TestXactService(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	ctx := context.TODO()

	// setup accounts and ledgers
	accountRepo := postgres.NewAccountRepository(db)
	johnDoe := account.Account{
		AccountID: account.AccountID("johndoe"),
		Currency:  currency.USD,
	}
	_, err = accountRepo.CreateAccount(ctx, johnDoe)
	require.NoError(t, err)
	maryJane := account.Account{
		AccountID: account.AccountID("maryjane"),
		Currency:  currency.USD,
	}
	_, err = accountRepo.CreateAccount(ctx, maryJane)
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
	balService := balance.NewService(balRepo)

	xactRepo := postgres.NewXactRepository(db)
	xactService := transaction.NewService(accountRepo, ledgerRepo, xactRepo, balRepo)

	t.Run("make deposit", func(t *testing.T) {
		xactNo := transaction.XactNo("dp100")
		err = xactService.MakeDeposit(ctx, transaction.DepositXact{
			XactNo:    xactNo,
			AccountID: johnDoe.AccountID,
			Amount:    decimal.NewFromInt(100),
		})
		require.NoError(t, err)

		// verify balance
		bal, err := balService.GetAccntBal(ctx, johnDoe.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(100).Equal(bal.CurrentBal))
	})

	t.Run("make withdrawal", func(t *testing.T) {
		xactNo := transaction.XactNo("wd100")
		err = xactService.MakeWithdrawal(ctx, transaction.WithdrawalXact{
			XactNo:    xactNo,
			AccountID: johnDoe.AccountID,
			Amount:    decimal.NewFromInt(25),
		})
		require.NoError(t, err)

		// verify balance
		bal, err := balService.GetAccntBal(ctx, johnDoe.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(75).Equal(bal.CurrentBal))

		t.Run("insufficient balance", func(t *testing.T) {
			xactNo := transaction.XactNo("wd102")
			err = xactService.MakeWithdrawal(ctx, transaction.WithdrawalXact{
				XactNo:    xactNo,
				AccountID: johnDoe.AccountID,
				Amount:    decimal.NewFromInt(200),
			})
			assert.ErrorIs(t, transaction.ErrInsufficientBalance, err)
		})
	})

	t.Run("make transfer", func(t *testing.T) {
		xactNo := transaction.XactNo("tr100")
		err = xactService.MakeTransfer(ctx, transaction.TransferXact{
			XactNo:    xactNo,
			FromAccnt: johnDoe.AccountID,
			ToAccnt:   maryJane.AccountID,
			Amount:    decimal.NewFromInt(25),
		})
		require.NoError(t, err)

		// verify john's balance
		johnBal, err := balService.GetAccntBal(ctx, johnDoe.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(50).Equal(johnBal.CurrentBal))

		// veryfy mary's balance
		maryBal, err := balService.GetAccntBal(ctx, maryJane.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(25).Equal(maryBal.CurrentBal))

		t.Run("insufficient balance", func(t *testing.T) {
			xactNo := transaction.XactNo("tr102")
			err = xactService.MakeWithdrawal(ctx, transaction.WithdrawalXact{
				XactNo:    xactNo,
				AccountID: johnDoe.AccountID,
				Amount:    decimal.NewFromInt(100),
			})
			assert.ErrorIs(t, transaction.ErrInsufficientBalance, err)
		})
	})
}
