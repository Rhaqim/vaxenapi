package providers

import (
	"log"

	"vaxen/api/internal/config"
	"vaxen/api/internal/providers/exchange"
	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/providers/mfa"
	"vaxen/api/internal/providers/payment"
	"vaxen/api/internal/providers/wallet"
)

// Bootstrap creates and configures a Registry from application config.
// Each provider is selected based on the corresponding env var.
// To swap a provider, change the env var and restart — or call Set* at runtime.
func Bootstrap(cfg *config.Config) *Registry {
	reg := NewRegistry()

	// --- KYC ---
	switch cfg.Providers.KYCProvider {
	case "sumsub":
		reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{
			APIKey:    cfg.Providers.SumSubAPIKey,
			APISecret: cfg.Providers.SumSubAPISecret,
			BaseURL:   cfg.Providers.SumSubBaseURL,
		}))
	default:
		log.Printf("Using default KYC provider: sumsub")
		reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	}

	// --- MFA ---
	switch cfg.Providers.MFAProvider {
	case "twilio":
		reg.SetMFA(mfa.NewTwilio(mfa.TwilioConfig{
			AccountSID: cfg.Providers.TwilioAccountSID,
			AuthToken:  cfg.Providers.TwilioAuthToken,
			ServiceSID: cfg.Providers.TwilioServiceSID,
		}))
	default:
		log.Printf("Using default MFA provider: twilio")
		reg.SetMFA(mfa.NewTwilio(mfa.TwilioConfig{}))
	}

	// --- Wallet ---
	switch cfg.Providers.WalletProvider {
	case "aws_kms":
		reg.SetWallet(wallet.NewAWSKMS(wallet.AWSKMSConfig{
			Region:    cfg.Providers.AWSRegion,
			KeyPolicy: cfg.Providers.AWSKMSKeyPolicy,
		}))
	default:
		log.Printf("Using default wallet provider: aws_kms")
		reg.SetWallet(wallet.NewAWSKMS(wallet.AWSKMSConfig{}))
	}

	// --- Exchange ---
	switch cfg.Providers.ExchangeProvider {
	case "internal":
		reg.SetExchange(exchange.NewInternal())
	default:
		log.Printf("Using default exchange provider: internal")
		reg.SetExchange(exchange.NewInternal())
	}

	// --- Payment ---
	switch cfg.Providers.PaymentProvider {
	case "circle":
		reg.SetPayment(payment.NewCircle(payment.CircleConfig{
			APIKey:  cfg.Providers.CircleAPIKey,
			BaseURL: cfg.Providers.CircleBaseURL,
		}))
	default:
		log.Printf("Using default payment provider: circle")
		reg.SetPayment(payment.NewCircle(payment.CircleConfig{}))
	}

	return reg
}
