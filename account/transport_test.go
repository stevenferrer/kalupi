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

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/etc/txdb"
	"github.com/sf9v/kalupi/postgres"
)

func TestHTTPHandler(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	err := postgres.Migrate(db)
	require.NoError(t, err)

	accountRepo := postgres.NewAccountRepository(db)
	accountService := account.NewService(accountRepo)

	logger := log.NewNopLogger()
	accountService = account.NewLoggingService(logger, accountService)

	accountHandler := account.NewHandler(accountService, logger)

	ctx := context.TODO()

	accountID := "johndoe"
	t.Run("create account", func(t *testing.T) {
		var req = map[string]interface{}{
			"accountId": accountID,
			"currency":  "USD",
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/accounts", bytes.NewBuffer(b))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		accountHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("get account", func(t *testing.T) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "/accounts/"+accountID, nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		accountHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp = struct {
			Account struct {
				AccountID string `json:"accountId"`
				Currency  string `json:"currency"`
			} `json:"account"`
			Err string `json:"error"`
		}{}
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Empty(t, resp.Err)

		require.NotNil(t, resp.Account)
		assert.Equal(t, accountID, resp.Account.AccountID)
		assert.Equal(t, "USD", resp.Account.Currency)
	})

	t.Run("list accounts", func(t *testing.T) {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "/accounts", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		accountHandler.ServeHTTP(rr, httpReq)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp = struct {
			Accounts []struct {
				AccountID string `json:"accountId"`
				Currency  string `json:"currency"`
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
		}
	})
}
