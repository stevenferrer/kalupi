package transaction_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
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

func TestHTTPHandler(t *testing.T) {
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
	ledgerService := ledger.NewService(ledgerRepo)
	err = ledgerService.CreateCashLedgers(ctx)
	require.NoError(t, err)

	balRepo := postgres.NewBalanceRepository(db)
	balService := balance.NewService(balRepo)

	xactRepo := postgres.NewXactRepository(db)
	xactService := transaction.NewService(accountRepo, ledgerRepo, xactRepo, balRepo)

	logger := log.NewNopLogger()
	xactService = transaction.NewLoggingService(logger, xactService)

	xactHandler := transaction.NewHandler(xactService, logger)

	t.Run("make deposit", func(t *testing.T) {
		var req = map[string]interface{}{
			"account_id": john.AccountID,
			"amount":     100,
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/deposit", bytes.NewBuffer(b))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		xactHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		// check balance
		johnBal, err := balService.GetAccntBal(ctx, john.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(100).Equal(johnBal.CurrentBal), "john should now have a balance of 100")

		t.Run("validation error", func(t *testing.T) {
			req = map[string]interface{}{
				"account_id": john.AccountID,
				"amount":     0,
			}

			b, err = json.Marshal(req)
			require.NoError(t, err)

			httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, "/deposit", bytes.NewBuffer(b))
			require.NoError(t, err)

			rr = httptest.NewRecorder()
			xactHandler.ServeHTTP(rr, httpReq)
			require.Equal(t, http.StatusBadRequest, rr.Code)
		})
	})

	t.Run("make withdrawal", func(t *testing.T) {
		var req = map[string]interface{}{
			"account_id": john.AccountID,
			"amount":     25,
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/withdraw", bytes.NewBuffer(b))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		xactHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		// check balance
		johnBal, err := balService.GetAccntBal(ctx, john.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(75).Equal(johnBal.CurrentBal), "john should now have a balance of 75")

		t.Run("validation error", func(t *testing.T) {
			req = map[string]interface{}{
				"account_id": john.AccountID,
				"amount":     0,
			}

			b, err = json.Marshal(req)
			require.NoError(t, err)

			httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, "/withdraw", bytes.NewBuffer(b))
			require.NoError(t, err)

			rr = httptest.NewRecorder()
			xactHandler.ServeHTTP(rr, httpReq)
			require.Equal(t, http.StatusBadRequest, rr.Code)
		})
	})

	t.Run("make payment", func(t *testing.T) {
		var req = map[string]interface{}{
			"from_account": john.AccountID,
			"to_account":   mary.AccountID,
			"amount":       30,
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/payment", bytes.NewBuffer(b))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		xactHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		// check balance
		johnBal, err := balService.GetAccntBal(ctx, john.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(45).Equal(johnBal.CurrentBal), "john should now have a balance of 100")

		maryBal, err := balService.GetAccntBal(ctx, mary.AccountID)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromInt(30).Equal(maryBal.CurrentBal), "mary should now have a balance of 30")

		t.Run("validation error", func(t *testing.T) {
			req = map[string]interface{}{
				"from_account": john.AccountID,
				"to_account":   mary.AccountID,
				"amount":       0,
			}

			b, err = json.Marshal(req)
			require.NoError(t, err)

			httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, "/payment", bytes.NewBuffer(b))
			require.NoError(t, err)

			rr = httptest.NewRecorder()
			xactHandler.ServeHTTP(rr, httpReq)
			require.Equal(t, http.StatusBadRequest, rr.Code)
		})
	})
}
