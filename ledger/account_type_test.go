package ledger_test

import (
	"testing"

	"github.com/sf9v/kalupi/ledger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountType(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tc := []struct {
			at     ledger.AccountType
			expect string
		}{
			{
				at:     ledger.AccountType(0),
				expect: "invalid",
			},
			{
				at:     ledger.AccountTypeLiability,
				expect: "AL",
			},
		}

		for _, tt := range tc {
			assert.Equal(t, tt.expect, tt.at.String())
		}
	})

	t.Run("sql scanner", func(t *testing.T) {
		var at ledger.AccountType
		err := at.Scan("AL")
		require.NoError(t, err)
		assert.Equal(t, ledger.AccountTypeLiability, at)
	})

	t.Run("driver valuer", func(t *testing.T) {
		value, err := ledger.AccountTypeLiability.Value()
		require.NoError(t, err)
		assert.NotNil(t, value)
	})
}
