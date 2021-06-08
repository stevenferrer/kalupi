package currency_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/kalupi/currency"
)

func TestCurrency(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tc := []struct {
			c      currency.Currency
			expect string
		}{
			{
				c:      currency.Currency(0),
				expect: "invalid",
			},
			{
				c:      currency.USD,
				expect: "USD",
			},
			// {
			// 	c:      currency.EUR,
			// 	expect: "EUR",
			// },
		}

		for _, tt := range tc {
			assert.Equal(t, tt.expect, tt.c.String())
		}
	})

	t.Run("is valid", func(t *testing.T) {
		tc := []struct {
			c      currency.Currency
			expect bool
		}{
			{
				c:      currency.Currency(0),
				expect: false,
			},
			{
				c:      currency.USD,
				expect: true,
			},
			// {
			// 	c:      currency.EUR,
			// 	expect: true,
			// },
		}

		for _, tt := range tc {
			assert.Equal(t, tt.expect, tt.c.IsValid())
		}
	})

	t.Run("sql scanner", func(t *testing.T) {
		var c currency.Currency
		err := c.Scan("USD")
		require.NoError(t, err)
		assert.Equal(t, currency.USD, c)
	})

	t.Run("driver valuer", func(t *testing.T) {
		value, err := currency.USD.Value()
		require.NoError(t, err)
		assert.NotNil(t, value)
	})

	t.Run("json", func(t *testing.T) {
		t.Run("marshal", func(t *testing.T) {
			b, err := json.Marshal(currency.USD)
			require.NoError(t, err)
			assert.Equal(t, `"USD"`, string(b))
		})

		t.Run("unmarshal", func(t *testing.T) {
			var c currency.Currency
			err := json.Unmarshal([]byte(`"USD"`), &c)
			require.NoError(t, err)
			assert.Equal(t, currency.USD, c)
		})
	})

}
