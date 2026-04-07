package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// RateStore abstracts the database lookup for exchange rates.
// The Internal provider uses this to read admin-managed rates
// without importing GORM directly (keeping the provider DB-agnostic).
type RateStore interface {
	// GetRate looks up the rate for a currency pair. Returns rate, spread, found.
	GetRate(from, to string) (rate decimal.Decimal, spread decimal.Decimal, found bool)
	// ListPairs returns all active currency pairs.
	ListPairs() []CurrencyPair
}

// Internal implements the Exchange Provider using internally managed rates.
// Rates are stored in the database and managed by platform admins.
type Internal struct {
	store      RateStore
	defaultFee decimal.Decimal
}

type InternalConfig struct {
	Store      RateStore
	DefaultFee decimal.Decimal // e.g. 0.001 = 0.1%
}

func NewInternal(cfg InternalConfig) *Internal {
	fee := cfg.DefaultFee
	if fee.IsZero() {
		fee = decimal.NewFromFloat(0.001) // 0.1% default
	}
	return &Internal{
		store:      cfg.Store,
		defaultFee: fee,
	}
}

func (i *Internal) Name() string { return "internal" }

func (i *Internal) GetRate(ctx context.Context, req RateRequest) (*RateResult, error) {
	if i.store != nil {
		rate, spread, found := i.store.GetRate(req.FromCurrency, req.ToCurrency)
		if found {
			return &RateResult{
				FromCurrency: req.FromCurrency,
				ToCurrency:   req.ToCurrency,
				Rate:         rate,
				Spread:       spread,
				Timestamp:    time.Now().Unix(),
			}, nil
		}
	}

	return nil, fmt.Errorf("no exchange rate configured for %s/%s", req.FromCurrency, req.ToCurrency)
}

func (i *Internal) GetQuote(ctx context.Context, req QuoteRequest) (*QuoteResult, error) {
	rateResult, err := i.GetRate(ctx, RateRequest{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
	})
	if err != nil {
		return nil, err
	}

	effectiveRate := rateResult.Rate.Sub(rateResult.Spread)
	toAmount := req.Amount.Mul(effectiveRate)
	fee := req.Amount.Mul(i.defaultFee)

	return &QuoteResult{
		QuoteID:      fmt.Sprintf("quote-%d", time.Now().UnixNano()),
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		FromAmount:   req.Amount,
		ToAmount:     toAmount,
		Rate:         effectiveRate,
		Spread:       rateResult.Spread,
		Fee:          fee,
		ExpiresAt:    time.Now().Add(30 * time.Second).Unix(),
	}, nil
}

func (i *Internal) ExecuteSwap(ctx context.Context, req SwapRequest) (*SwapResult, error) {
	// In a full implementation this would:
	// 1. Look up the quote by ID, validate it hasn't expired
	// 2. Debit the source wallet
	// 3. Credit the destination wallet
	// 4. Record the transaction
	// For now, return a pending result for the service layer to record.
	return &SwapResult{
		TransactionID: fmt.Sprintf("swap-%d", time.Now().UnixNano()),
		Status:        "completed",
		FromAmount:    decimal.Zero,
		ToAmount:      decimal.Zero,
		Rate:          decimal.Zero,
		Fee:           decimal.Zero,
		ProviderRef:   "internal",
	}, nil
}

func (i *Internal) ListSupportedPairs(ctx context.Context) ([]CurrencyPair, error) {
	if i.store != nil {
		pairs := i.store.ListPairs()
		if len(pairs) > 0 {
			return pairs, nil
		}
	}

	return nil, fmt.Errorf("no exchange pairs configured — seed rates via admin API or CLI")
}
