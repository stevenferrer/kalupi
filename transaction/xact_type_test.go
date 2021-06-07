package transaction_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/transaction"
)

func TestXactType(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tc := []struct {
			tt     transaction.XactType
			expect string
		}{
			{
				tt:     transaction.XactType(0),
				expect: "invalid",
			},
			{
				tt:     transaction.XactTypeDebit,
				expect: "Dr",
			},
			{
				tt:     transaction.XactTypeCredit,
				expect: "Cr",
			},
		}

		for _, tt := range tc {
			assert.Equal(t, tt.expect, tt.tt.String())
		}
	})

	t.Run("sql scanner", func(t *testing.T) {
		tc := []struct {
			s      string
			expect transaction.XactType
		}{
			{
				s:      "Dr",
				expect: transaction.XactTypeDebit,
			},
			{
				s:      "Cr",
				expect: transaction.XactTypeCredit,
			},
		}

		for _, tt := range tc {
			var ttyp transaction.XactType
			err := ttyp.Scan(tt.s)
			require.NoError(t, err)
			assert.Equal(t, tt.expect, ttyp)
		}
	})

	t.Run("driver valuer", func(t *testing.T) {
		value, err := transaction.XactTypeDebit.Value()
		require.NoError(t, err)
		assert.NotNil(t, value)
	})
}
