package services_test

import (
	"context"
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/exchange"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupExchangeService(t *testing.T) *services.ExchangeService {
	t.Helper()
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.ConversionOrder{})
	reg := providers.NewRegistry()
	reg.SetExchange(exchange.NewInternal())
	return services.NewExchangeService(db, reg)
}

func TestExchangeService_GetRate(t *testing.T) {
	svc := setupExchangeService(t)

	rate, err := svc.GetRate(context.Background(), "USD", "EUR")
	require.NoError(t, err)
	assert.Equal(t, "USD", rate.FromCurrency)
	assert.Equal(t, "EUR", rate.ToCurrency)
	assert.True(t, rate.Rate.GreaterThan(decimal.Zero))
}

func TestExchangeService_GetQuote(t *testing.T) {
	svc := setupExchangeService(t)

	quote, err := svc.GetQuote(context.Background(), "USD", "EUR", decimal.NewFromInt(1000), "sell")
	require.NoError(t, err)
	assert.NotEmpty(t, quote.QuoteID)
	assert.True(t, quote.ToAmount.GreaterThan(decimal.Zero))
	assert.True(t, quote.Fee.GreaterThan(decimal.Zero))
}

func TestExchangeService_ExecuteSwap(t *testing.T) {
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.ConversionOrder{})
	reg := providers.NewRegistry()
	reg.SetExchange(exchange.NewInternal())
	svc := services.NewExchangeService(db, reg)
	org := testutil.SeedOrganization(t, db)

	order, err := svc.ExecuteSwap(context.Background(), org.ID, "quote-123")
	require.NoError(t, err)
	assert.NotEmpty(t, order.ID)
}

func TestExchangeService_ListSupportedPairs(t *testing.T) {
	svc := setupExchangeService(t)

	pairs, err := svc.ListSupportedPairs(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, pairs)
}

func TestExchangeService_NilProvider(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	svc := services.NewExchangeService(db, reg)

	_, err := svc.GetRate(context.Background(), "USD", "EUR")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}
