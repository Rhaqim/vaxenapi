package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Circle implements the Payment Provider interface using Circle's API.
// See: https://developers.circle.com/
type Circle struct {
	apiKey  string
	baseURL string
}

type CircleConfig struct {
	APIKey  string
	BaseURL string // e.g. "https://api.circle.com" or sandbox
}

func NewCircle(cfg CircleConfig) *Circle {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.circle.com"
	}
	return &Circle{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
	}
}

func (c *Circle) Name() string { return "circle" }

func (c *Circle) SendPayment(ctx context.Context, req SendRequest) (*SendResult, error) {
	// TODO: Implement Circle Payouts API
	// POST /v1/payouts
	return &SendResult{
		ProviderRef:   fmt.Sprintf("circle-payout-%d", time.Now().UnixNano()),
		Status:        "pending",
		Fee:           decimal.NewFromFloat(0.25),
		EstimatedTime: "1-2 business days",
	}, nil
}

func (c *Circle) GetPaymentStatus(ctx context.Context, providerRef string) (*PaymentStatus, error) {
	// TODO: Implement Circle Payouts status check
	// GET /v1/payouts/{id}
	return &PaymentStatus{
		ProviderRef: providerRef,
		Status:      "pending",
		UpdatedAt:   time.Now().Unix(),
	}, nil
}

func (c *Circle) ValidateBeneficiary(ctx context.Context, req ValidateBeneficiaryRequest) (*ValidationResult, error) {
	// TODO: Implement Circle bank account validation
	return &ValidationResult{
		Valid:   true,
		Message: "Account details are valid",
	}, nil
}

func (c *Circle) GetFees(ctx context.Context, req FeeRequest) (*FeeResult, error) {
	// TODO: Implement Circle fee estimation
	return &FeeResult{
		Fee:      decimal.NewFromFloat(0.25),
		Currency: req.Currency,
		Method:   req.Method,
	}, nil
}

func (c *Circle) CancelPayment(ctx context.Context, providerRef string) error {
	// TODO: Implement Circle payout cancellation
	// POST /v1/payouts/{id}/cancel
	return nil
}
