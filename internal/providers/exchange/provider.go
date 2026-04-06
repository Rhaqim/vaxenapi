package exchange

import (
	"context"

	"github.com/shopspring/decimal"
)

// Provider defines the interface for currency exchange / swap services.
// Implementations can wrap internal rates, Binance, Kraken, CurrencyCloud, etc.
type Provider interface {
	Name() string

	// GetRate returns the current exchange rate between two currencies.
	GetRate(ctx context.Context, req RateRequest) (*RateResult, error)

	// GetQuote returns a locked-in quote for a specific amount and pair.
	GetQuote(ctx context.Context, req QuoteRequest) (*QuoteResult, error)

	// ExecuteSwap executes a currency swap using a previously obtained quote.
	ExecuteSwap(ctx context.Context, req SwapRequest) (*SwapResult, error)

	// ListSupportedPairs returns all tradable currency pairs.
	ListSupportedPairs(ctx context.Context) ([]CurrencyPair, error)
}

type RateRequest struct {
	FromCurrency string
	ToCurrency   string
}

type RateResult struct {
	FromCurrency string
	ToCurrency   string
	Rate         decimal.Decimal
	Spread       decimal.Decimal
	Timestamp    int64
}

type QuoteRequest struct {
	FromCurrency string
	ToCurrency   string
	Amount       decimal.Decimal
	Side         string // "buy" or "sell"
}

type QuoteResult struct {
	QuoteID      string
	FromCurrency string
	ToCurrency   string
	FromAmount   decimal.Decimal
	ToAmount     decimal.Decimal
	Rate         decimal.Decimal
	Spread       decimal.Decimal
	Fee          decimal.Decimal
	ExpiresAt    int64 // unix timestamp
}

type SwapRequest struct {
	QuoteID        string
	OrganizationID string
}

type SwapResult struct {
	TransactionID string
	Status        string // "pending", "completed", "failed"
	FromAmount    decimal.Decimal
	ToAmount      decimal.Decimal
	Rate          decimal.Decimal
	Fee           decimal.Decimal
	ProviderRef   string
}

type CurrencyPair struct {
	From     string
	To       string
	MinAmount decimal.Decimal
	MaxAmount decimal.Decimal
	Type      string // "fiat", "crypto", "cross"
}
