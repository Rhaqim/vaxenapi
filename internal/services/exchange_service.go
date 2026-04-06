package services

import (
	"context"
	"errors"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/exchange"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ExchangeService struct {
	db  *gorm.DB
	reg *providers.Registry
}

func NewExchangeService(db *gorm.DB, reg *providers.Registry) *ExchangeService {
	return &ExchangeService{db: db, reg: reg}
}

// GetRate returns the current exchange rate between two currencies.
func (s *ExchangeService) GetRate(ctx context.Context, from, to string) (*exchange.RateResult, error) {
	provider := s.reg.Exchange()
	if provider == nil {
		return nil, errors.New("exchange provider not configured")
	}

	return provider.GetRate(ctx, exchange.RateRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
}

// GetQuote returns a locked-in quote for a conversion.
func (s *ExchangeService) GetQuote(ctx context.Context, from, to string, amount decimal.Decimal, side string) (*exchange.QuoteResult, error) {
	provider := s.reg.Exchange()
	if provider == nil {
		return nil, errors.New("exchange provider not configured")
	}

	return provider.GetQuote(ctx, exchange.QuoteRequest{
		FromCurrency: from,
		ToCurrency:   to,
		Amount:       amount,
		Side:         side,
	})
}

// ExecuteSwap executes a swap and records it as a ConversionOrder.
func (s *ExchangeService) ExecuteSwap(ctx context.Context, orgID string, quoteID string) (*models.ConversionOrder, error) {
	provider := s.reg.Exchange()
	if provider == nil {
		return nil, errors.New("exchange provider not configured")
	}

	result, err := provider.ExecuteSwap(ctx, exchange.SwapRequest{
		QuoteID:        quoteID,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, err
	}

	order := models.ConversionOrder{
		OrganizationID: orgID,
		FromCurrency:   "", // populated from quote
		ToCurrency:     "",
		FromAmount:     result.FromAmount,
		ToAmount:       result.ToAmount,
		Rate:           result.Rate,
		Spread:         decimal.Zero,
		Fee:            result.Fee,
		Status:         models.OrderStatusCompleted,
		ProviderRef:    &result.ProviderRef,
	}
	if err := s.db.Create(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// ListSupportedPairs returns all tradable currency pairs.
func (s *ExchangeService) ListSupportedPairs(ctx context.Context) ([]exchange.CurrencyPair, error) {
	provider := s.reg.Exchange()
	if provider == nil {
		return nil, errors.New("exchange provider not configured")
	}

	return provider.ListSupportedPairs(ctx)
}
