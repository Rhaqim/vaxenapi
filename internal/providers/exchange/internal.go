package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Internal implements the Exchange Provider using internally managed rates.
// Rates can be set by platform admins. This is the default provider
// that can be replaced with Binance, Kraken, CurrencyCloud, etc.
type Internal struct {
	// In production, rates would be stored in DB and managed by admin.
	// This stub uses hardcoded rates for demonstration.
}

func NewInternal() *Internal {
	return &Internal{}
}

func (i *Internal) Name() string { return "internal" }

func (i *Internal) GetRate(ctx context.Context, req RateRequest) (*RateResult, error) {
	// TODO: Look up rate from database (managed by platform admin)
	rate := decimal.NewFromFloat(1.0)
	spread := decimal.NewFromFloat(0.005) // 0.5% spread

	return &RateResult{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
		Spread:       spread,
		Timestamp:    time.Now().Unix(),
	}, nil
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
	fee := req.Amount.Mul(decimal.NewFromFloat(0.001)) // 0.1% fee

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
	// TODO: Execute the swap using the stored quote
	// 1. Validate quote hasn't expired
	// 2. Debit source wallet
	// 3. Credit destination wallet
	// 4. Record the transaction
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
	// TODO: Load from database
	return []CurrencyPair{
		{From: "USD", To: "EUR", MinAmount: decimal.NewFromInt(1), MaxAmount: decimal.NewFromInt(1000000), Type: "fiat"},
		{From: "USD", To: "GBP", MinAmount: decimal.NewFromInt(1), MaxAmount: decimal.NewFromInt(1000000), Type: "fiat"},
		{From: "USD", To: "ETH", MinAmount: decimal.NewFromInt(10), MaxAmount: decimal.NewFromInt(500000), Type: "cross"},
		{From: "USD", To: "BTC", MinAmount: decimal.NewFromInt(10), MaxAmount: decimal.NewFromInt(500000), Type: "cross"},
		{From: "ETH", To: "BTC", MinAmount: decimal.NewFromFloat(0.01), MaxAmount: decimal.NewFromInt(1000), Type: "crypto"},
	}, nil
}
