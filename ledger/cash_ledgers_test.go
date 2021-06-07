package ledger_test

import (
	"testing"

	"github.com/sf9v/kalupi/currency"
	"github.com/sf9v/kalupi/ledger"
	"github.com/stretchr/testify/assert"
)

func TestGetCashLedgerNo(t *testing.T) {
	tc := []struct {
		curr     currency.Currency
		expect   ledger.LedgerNo
		hasError bool
	}{
		{
			curr:     currency.Currency(0),
			hasError: true,
		},
		{
			curr:   currency.USD,
			expect: ledger.CashUSDLedgerNo,
		},
	}

	for _, tt := range tc {
		got, err := ledger.GetCashLedgerNo(tt.curr)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, got)
		}
	}
}
