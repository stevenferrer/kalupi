package balance_test

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

func TestBalanceService(t *testing.T) {
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
	xactService := transaction.NewService(accountRepo,
		ledgerRepo, xactRepo, balRepo)

	// john deposits 100USD
	err = xactService.MakeDeposit(ctx, transaction.DepositXact{
		AccountID: johnDoe.AccountID,
		Amount:    decimal.NewFromInt(100),
	})
	require.NoError(t, err)

	// john withdraws 25USD
	err = xactService.MakeWithdrawal(ctx, transaction.WithdrawalXact{
		AccountID: johnDoe.AccountID,
		Amount:    decimal.NewFromInt(25),
	})
	require.NoError(t, err)

	// john sends 25USD to mary
	err = xactService.MakeTransfer(ctx, transaction.TransferXact{
		FromAccount: johnDoe.AccountID,
		ToAccount:   maryJane.AccountID,
		Amount:      decimal.NewFromInt(25),
	})
	require.NoError(t, err)

	// john and mary should have 50USD and 25USD, respectively
	johnBal, err := balService.GetAccntBal(ctx, johnDoe.AccountID)
	require.NoError(t, err)
	assert.True(t, decimal.NewFromInt(50).Equal(johnBal.CurrentBal))

	maryBal, err := balService.GetAccntBal(ctx, maryJane.AccountID)
	require.NoError(t, err)
	assert.True(t, decimal.NewFromInt(25).Equal(maryBal.CurrentBal))
}
