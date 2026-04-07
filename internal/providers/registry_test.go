package providers_test

import (
	"testing"

	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/email"
	"vaxen/api/internal/providers/exchange"
	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/providers/mfa"
	"vaxen/api/internal/providers/payment"
	"vaxen/api/internal/providers/wallet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Empty_Validate(t *testing.T) {
	reg := providers.NewRegistry()
	err := reg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "KYC provider not configured")
}

func TestRegistry_Partial_Validate(t *testing.T) {
	reg := providers.NewRegistry()
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	err := reg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MFA provider not configured")
}

func TestRegistry_Full_Validate(t *testing.T) {
	reg := fullyConfiguredRegistry()
	err := reg.Validate()
	assert.NoError(t, err)
}

func TestRegistry_SetAndGet(t *testing.T) {
	reg := providers.NewRegistry()

	// Initially nil
	assert.Nil(t, reg.KYC())
	assert.Nil(t, reg.MFA())
	assert.Nil(t, reg.Wallet())
	assert.Nil(t, reg.Exchange())
	assert.Nil(t, reg.Payment())
	assert.Nil(t, reg.Email())

	// Set providers
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	reg.SetMFA(mfa.NewTwilio(mfa.TwilioConfig{}))
	reg.SetWallet(wallet.NewAWSKMS(wallet.AWSKMSConfig{}))
	reg.SetExchange(exchange.NewInternal(exchange.InternalConfig{}))
	reg.SetPayment(payment.NewCircle(payment.CircleConfig{}))
	reg.SetEmail(email.NewSendGrid(email.SendGridConfig{}))

	// Now all set
	assert.NotNil(t, reg.KYC())
	assert.NotNil(t, reg.MFA())
	assert.NotNil(t, reg.Wallet())
	assert.NotNil(t, reg.Exchange())
	assert.NotNil(t, reg.Payment())
	assert.NotNil(t, reg.Email())

	assert.Equal(t, "sumsub", reg.KYC().Name())
	assert.Equal(t, "twilio", reg.MFA().Name())
	assert.Equal(t, "aws_kms", reg.Wallet().Name())
	assert.Equal(t, "internal", reg.Exchange().Name())
	assert.Equal(t, "circle", reg.Payment().Name())
	assert.Equal(t, "sendgrid", reg.Email().Name())
}

func TestRegistry_Swap(t *testing.T) {
	reg := fullyConfiguredRegistry()
	require.Equal(t, "sumsub", reg.KYC().Name())

	// Swap KYC provider at runtime
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{APIKey: "new-key"}))
	assert.Equal(t, "sumsub", reg.KYC().Name())
}

func fullyConfiguredRegistry() *providers.Registry {
	reg := providers.NewRegistry()
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	reg.SetMFA(mfa.NewTwilio(mfa.TwilioConfig{}))
	reg.SetWallet(wallet.NewAWSKMS(wallet.AWSKMSConfig{}))
	reg.SetExchange(exchange.NewInternal(exchange.InternalConfig{}))
	reg.SetPayment(payment.NewCircle(payment.CircleConfig{}))
	reg.SetEmail(email.NewSendGrid(email.SendGridConfig{}))
	return reg
}
