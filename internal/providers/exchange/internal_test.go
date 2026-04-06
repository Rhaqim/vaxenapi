package exchange_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/exchange"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newInternal() *exchange.Internal {
	return exchange.NewInternal()
}

func TestInternal_Name(t *testing.T) {
	assert.Equal(t, "internal", newInternal().Name())
}

func TestInternal_ImplementsProvider(t *testing.T) {
	var _ exchange.Provider = newInternal()
}

func TestInternal_GetRate(t *testing.T) {
	ex := newInternal()
	result, err := ex.GetRate(context.Background(), exchange.RateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	})
	require.NoError(t, err)
	assert.Equal(t, "USD", result.FromCurrency)
	assert.Equal(t, "EUR", result.ToCurrency)
	assert.True(t, result.Rate.GreaterThan(decimal.Zero))
	assert.True(t, result.Spread.GreaterThan(decimal.Zero))
	assert.Greater(t, result.Timestamp, int64(0))
}

func TestInternal_GetQuote(t *testing.T) {
	ex := newInternal()
	amount := decimal.NewFromInt(1000)
	result, err := ex.GetQuote(context.Background(), exchange.QuoteRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       amount,
		Side:         "sell",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.QuoteID)
	assert.Equal(t, "USD", result.FromCurrency)
	assert.Equal(t, "EUR", result.ToCurrency)
	assert.True(t, result.FromAmount.Equal(amount))
	assert.True(t, result.ToAmount.GreaterThan(decimal.Zero))
	assert.True(t, result.Fee.GreaterThan(decimal.Zero))
	assert.Greater(t, result.ExpiresAt, int64(0))
}

func TestInternal_ExecuteSwap(t *testing.T) {
	ex := newInternal()
	result, err := ex.ExecuteSwap(context.Background(), exchange.SwapRequest{
		QuoteID:        "quote-123",
		OrganizationID: "org-456",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.TransactionID)
	assert.Equal(t, "completed", result.Status)
}

func TestInternal_ListSupportedPairs(t *testing.T) {
	ex := newInternal()
	pairs, err := ex.ListSupportedPairs(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, pairs)

	// Check we have all expected pair types
	types := make(map[string]bool)
	for _, p := range pairs {
		types[p.Type] = true
		assert.NotEmpty(t, p.From)
		assert.NotEmpty(t, p.To)
	}
	assert.True(t, types["fiat"], "should have fiat pairs")
	assert.True(t, types["crypto"], "should have crypto pairs")
	assert.True(t, types["cross"], "should have cross pairs")
}
