package payment_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/payment"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCircle() *payment.Circle {
	return payment.NewCircle(payment.CircleConfig{
		APIKey: "test-key",
	})
}

func TestCircle_Name(t *testing.T) {
	assert.Equal(t, "circle", newCircle().Name())
}

func TestCircle_ImplementsProvider(t *testing.T) {
	var _ payment.Provider = newCircle()
}

func TestCircle_SendPayment(t *testing.T) {
	c := newCircle()
	result, err := c.SendPayment(context.Background(), payment.SendRequest{
		OrganizationID:  "org-123",
		Amount:          decimal.NewFromFloat(500.00),
		Currency:        "USD",
		Reference:       "PAY-001",
		BeneficiaryName: "John Doe",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ProviderRef)
	assert.Equal(t, "pending", result.Status)
	assert.True(t, result.Fee.GreaterThan(decimal.Zero))
	assert.NotEmpty(t, result.EstimatedTime)
}

func TestCircle_GetPaymentStatus(t *testing.T) {
	c := newCircle()
	result, err := c.GetPaymentStatus(context.Background(), "ref-123")
	require.NoError(t, err)
	assert.Equal(t, "ref-123", result.ProviderRef)
	assert.Equal(t, "pending", result.Status)
	assert.Greater(t, result.UpdatedAt, int64(0))
}

func TestCircle_ValidateBeneficiary(t *testing.T) {
	c := newCircle()
	result, err := c.ValidateBeneficiary(context.Background(), payment.ValidateBeneficiaryRequest{
		IBAN:     "GB29NWBK60161331926819",
		Country:  "GB",
		Currency: "GBP",
	})
	require.NoError(t, err)
	assert.True(t, result.Valid)
}

func TestCircle_GetFees(t *testing.T) {
	c := newCircle()
	result, err := c.GetFees(context.Background(), payment.FeeRequest{
		Amount:   decimal.NewFromFloat(1000.00),
		Currency: "USD",
		Country:  "US",
		Method:   "wire",
	})
	require.NoError(t, err)
	assert.True(t, result.Fee.GreaterThan(decimal.Zero))
	assert.Equal(t, "USD", result.Currency)
	assert.Equal(t, "wire", result.Method)
}

func TestCircle_CancelPayment(t *testing.T) {
	c := newCircle()
	err := c.CancelPayment(context.Background(), "ref-123")
	assert.NoError(t, err)
}

func TestCircle_DefaultBaseURL(t *testing.T) {
	c := payment.NewCircle(payment.CircleConfig{})
	assert.Equal(t, "circle", c.Name())
}
