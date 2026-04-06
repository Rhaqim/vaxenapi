package payment

import (
	"context"

	"github.com/shopspring/decimal"
)

// Provider defines the interface for money transfer / payout services.
// Implementations can wrap Circle, Wise, Stripe Connect, CurrencyCloud, etc.
type Provider interface {
	Name() string

	// SendPayment initiates a money transfer to a beneficiary.
	SendPayment(ctx context.Context, req SendRequest) (*SendResult, error)

	// GetPaymentStatus checks the status of a previously initiated payment.
	GetPaymentStatus(ctx context.Context, providerRef string) (*PaymentStatus, error)

	// ValidateBeneficiary verifies that beneficiary details are valid before sending.
	ValidateBeneficiary(ctx context.Context, req ValidateBeneficiaryRequest) (*ValidationResult, error)

	// GetFees returns the estimated fee for a payment.
	GetFees(ctx context.Context, req FeeRequest) (*FeeResult, error)

	// CancelPayment attempts to cancel a pending payment.
	CancelPayment(ctx context.Context, providerRef string) error
}

type SendRequest struct {
	OrganizationID string
	Amount         decimal.Decimal
	Currency       string
	Reference      string
	Description    string

	// Beneficiary details
	BeneficiaryID   string
	BeneficiaryName string
	AccountNumber   string
	RoutingNumber   string
	BankCode        string
	IBAN            string
	SwiftBIC        string
	Country         string

	// Crypto-specific
	WalletAddress string
	Network       string
}

type SendResult struct {
	ProviderRef   string
	Status        string // "pending", "processing", "completed", "failed"
	Fee           decimal.Decimal
	EstimatedTime string // e.g. "1-2 business days"
}

type PaymentStatus struct {
	ProviderRef string
	Status      string
	UpdatedAt   int64
	FailReason  string
}

type ValidateBeneficiaryRequest struct {
	AccountNumber string
	RoutingNumber string
	IBAN          string
	SwiftBIC      string
	Country       string
	Currency      string
}

type ValidationResult struct {
	Valid   bool
	Message string
}

type FeeRequest struct {
	Amount   decimal.Decimal
	Currency string
	Country  string
	Method   string // "wire", "ach", "sepa", "crypto"
}

type FeeResult struct {
	Fee      decimal.Decimal
	Currency string
	Method   string
}
