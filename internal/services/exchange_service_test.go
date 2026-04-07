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
	"gorm.io/gorm"
)

func setupExchangeService(t *testing.T) (*services.ExchangeService, *gorm.DB) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.ConversionOrder{}, &models.ExchangeRate{})

	// Seed exchange rates
	rates := []models.ExchangeRate{
		{FromCurrency: "USD", ToCurrency: "EUR", Rate: decimal.NewFromFloat(0.85), Spread: decimal.NewFromFloat(0.005), IsActive: true},
		{FromCurrency: "USD", ToCurrency: "GBP", Rate: decimal.NewFromFloat(0.73), Spread: decimal.NewFromFloat(0.004), IsActive: true},
		{FromCurrency: "USD", ToCurrency: "ETH", Rate: decimal.NewFromFloat(0.00033), Spread: decimal.NewFromFloat(0.001), IsActive: true},
	}
	for _, r := range rates {
		db.Create(&r)
	}

	store := exchange.NewGormRateStore(db)
	reg := providers.NewRegistry()
	reg.SetExchange(exchange.NewInternal(exchange.InternalConfig{Store: store}))
	return services.NewExchangeService(db, reg), db
}

func TestExchangeService_GetRate(t *testing.T) {
	svc, _ := setupExchangeService(t)

	rate, err := svc.GetRate(context.Background(), "USD", "EUR")
	require.NoError(t, err)
	assert.Equal(t, "USD", rate.FromCurrency)
	assert.Equal(t, "EUR", rate.ToCurrency)
	assert.True(t, rate.Rate.Equal(decimal.NewFromFloat(0.85)))
}

func TestExchangeService_GetRate_NotFound(t *testing.T) {
	svc, _ := setupExchangeService(t)

	_, err := svc.GetRate(context.Background(), "USD", "JPY")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no exchange rate configured")
}

func TestExchangeService_GetQuote(t *testing.T) {
	svc, _ := setupExchangeService(t)

	quote, err := svc.GetQuote(context.Background(), "USD", "EUR", decimal.NewFromInt(1000), "sell")
	require.NoError(t, err)
	assert.NotEmpty(t, quote.QuoteID)
	assert.True(t, quote.ToAmount.GreaterThan(decimal.Zero))
	assert.True(t, quote.Fee.GreaterThan(decimal.Zero))
}

func TestExchangeService_ExecuteSwap(t *testing.T) {
	svc, db := setupExchangeService(t)
	org := testutil.SeedOrganization(t, db)

	order, err := svc.ExecuteSwap(context.Background(), org.ID, "quote-123")
	require.NoError(t, err)
	assert.NotEmpty(t, order.ID)
}

func TestExchangeService_ListSupportedPairs(t *testing.T) {
	svc, _ := setupExchangeService(t)

	pairs, err := svc.ListSupportedPairs(context.Background())
	require.NoError(t, err)
	assert.Len(t, pairs, 3)
}

func TestExchangeService_NilProvider(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	svc := services.NewExchangeService(db, reg)

	_, err := svc.GetRate(context.Background(), "USD", "EUR")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}
