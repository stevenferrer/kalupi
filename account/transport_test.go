package account_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/kalupi/account"
	accountsvc "github.com/stevenferrer/kalupi/account/service"
	"github.com/stevenferrer/kalupi/balance"
	"github.com/stevenferrer/kalupi/etc/txdb"
	"github.com/stevenferrer/kalupi/postgres"
)

func TestHTTPHandler(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	balRepo := postgres.NewBalanceRepository(db)
	balService := balance.NewService(balRepo)

	accountRepo := postgres.NewAccountRepository(db)
	accountService := accountsvc.New(accountRepo, balService)

	logger := log.NewNopLogger()
	accountService = account.NewLoggingService(logger, accountService)

	accountHandler := account.NewHTTPHandler(accountService, logger)

	ctx := context.TODO()

	accountID := "johndoe"
	t.Run("create account", func(t *testing.T) {
		var req = map[string]interface{}{
			"account_id": accountID,
			"currency":   "USD",
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewBuffer(b))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		accountHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		t.Run("validation error", func(t *testing.T) {
			var req = map[string]interface{}{}
			b, err = json.Marshal(req)
			require.NoError(t, err)

			httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewBuffer(b))
			require.NoError(t, err)

			rr = httptest.NewRecorder()
			accountHandler.ServeHTTP(rr, httpReq)
			require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

			var resp = map[string]interface{}{}
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)
		})
	})

	t.Run("get account", func(t *testing.T) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "/"+accountID, nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		accountHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp = struct {
			Account struct {
				AccountID string `json:"id"`
				Currency  string `json:"currency"`
				Balance   string `json:"balance"`
			} `json:"account"`
			Err string `json:"error"`
		}{}
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Empty(t, resp.Err)

		require.NotNil(t, resp.Account)
		assert.Equal(t, accountID, resp.Account.AccountID)
		assert.Equal(t, "USD", resp.Account.Currency)
		assert.Equal(t, "0", resp.Account.Balance)

		t.Run("not found", func(t *testing.T) {
			httpReq, err = http.NewRequestWithContext(ctx, http.MethodGet, "/johntravolta", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			accountHandler.ServeHTTP(rr, httpReq)
			require.Equal(t, http.StatusNotFound, rr.Code)
		})
	})

	t.Run("list accounts", func(t *testing.T) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		accountHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp = struct {
			Accounts []struct {
				AccountID string `json:"id"`
				Currency  string `json:"currency"`
				Balance   string `json:"balance"`
			} `json:"accounts"`
			Err string `json:"error"`
		}{}
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Empty(t, resp.Err)

		require.Len(t, resp.Accounts, 1)
		for _, accnt := range resp.Accounts {
			assert.NotEmpty(t, accnt.AccountID)
			assert.NotEmpty(t, accnt.Currency)
			assert.Equal(t, "0", accnt.Balance)
		}
	})
}
