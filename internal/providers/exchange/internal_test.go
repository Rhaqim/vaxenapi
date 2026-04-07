package exchange_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/exchange"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRateStore implements exchange.RateStore for testing.
type mockRateStore struct {
	rates map[string]mockRate
}

type mockRate struct {
	rate   decimal.Decimal
	spread decimal.Decimal
}

func newMockStore() *mockRateStore {
	return &mockRateStore{
		rates: map[string]mockRate{
			"USD/EUR": {rate: decimal.NewFromFloat(0.85), spread: decimal.NewFromFloat(0.005)},
			"USD/GBP": {rate: decimal.NewFromFloat(0.73), spread: decimal.NewFromFloat(0.004)},
			"USD/ETH": {rate: decimal.NewFromFloat(0.00033), spread: decimal.NewFromFloat(0.001)},
			"ETH/BTC": {rate: decimal.NewFromFloat(0.055), spread: decimal.NewFromFloat(0.002)},
		},
	}
}

func (m *mockRateStore) GetRate(from, to string) (decimal.Decimal, decimal.Decimal, bool) {
	key := from + "/" + to
	r, ok := m.rates[key]
	if !ok {
		return decimal.Zero, decimal.Zero, false
	}
	return r.rate, r.spread, true
}

func (m *mockRateStore) ListPairs() []exchange.CurrencyPair {
	pairs := make([]exchange.CurrencyPair, 0)
	for key := range m.rates {
		parts := []byte(key)
		from := string(parts[:3])
		to := string(parts[4:])
		pairs = append(pairs, exchange.CurrencyPair{
			From: from, To: to,
			MinAmount: decimal.NewFromInt(1),
			MaxAmount: decimal.NewFromInt(1000000),
			Type:      "fiat",
		})
	}
	return pairs
}

func newInternal() *exchange.Internal {
	return exchange.NewInternal(exchange.InternalConfig{Store: newMockStore()})
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
	assert.True(t, result.Rate.Equal(decimal.NewFromFloat(0.85)))
	assert.True(t, result.Spread.Equal(decimal.NewFromFloat(0.005)))
	assert.Greater(t, result.Timestamp, int64(0))
}

func TestInternal_GetRate_NotFound(t *testing.T) {
	ex := newInternal()
	_, err := ex.GetRate(context.Background(), exchange.RateRequest{
		FromCurrency: "USD",
		ToCurrency:   "JPY",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no exchange rate configured")
}

func TestInternal_GetRate_NoStore(t *testing.T) {
	ex := exchange.NewInternal(exchange.InternalConfig{})
	_, err := ex.GetRate(context.Background(), exchange.RateRequest{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no exchange rate configured")
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
	assert.Len(t, pairs, 4)
}

func TestInternal_ListSupportedPairs_NoStore(t *testing.T) {
	ex := exchange.NewInternal(exchange.InternalConfig{})
	_, err := ex.ListSupportedPairs(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no exchange pairs configured")
}
