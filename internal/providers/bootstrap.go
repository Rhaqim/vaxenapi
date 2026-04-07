package providers

import (
	"log"

	"vaxen/api/internal/config"
	"vaxen/api/internal/providers/email"
	"vaxen/api/internal/providers/exchange"
	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/providers/mfa"
	"vaxen/api/internal/providers/payment"
	"vaxen/api/internal/providers/wallet"

	"gorm.io/gorm"
)

// Bootstrap creates and configures a Registry from application config.
// Each provider is selected based on the corresponding env var.
// To swap a provider, change the env var and restart — or call Set* at runtime.
func Bootstrap(cfg *config.Config, db *gorm.DB) *Registry {
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
	var rateStore exchange.RateStore
	if db != nil {
		rateStore = exchange.NewGormRateStore(db)
	}
	switch cfg.Providers.ExchangeProvider {
	case "internal":
		reg.SetExchange(exchange.NewInternal(exchange.InternalConfig{Store: rateStore}))
	default:
		log.Printf("Using default exchange provider: internal")
		reg.SetExchange(exchange.NewInternal(exchange.InternalConfig{Store: rateStore}))
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

	// --- Email ---
	switch cfg.Providers.EmailProvider {
	case "smtp":
		reg.SetEmail(email.NewSMTP(email.SMTPConfig{
			Host:      cfg.Providers.SMTPHost,
			Port:      cfg.Providers.SMTPPort,
			Username:  cfg.Providers.SMTPUsername,
			Password:  cfg.Providers.SMTPPassword,
			FromEmail: cfg.Providers.EmailFromAddress,
			FromName:  cfg.Providers.EmailFromName,
		}))
	default:
		log.Printf("Using default email provider: sendgrid")
		reg.SetEmail(email.NewSendGrid(email.SendGridConfig{
			APIKey:    cfg.Providers.SendGridAPIKey,
			FromEmail: cfg.Providers.EmailFromAddress,
			FromName:  cfg.Providers.EmailFromName,
		}))
	}

	return reg
}
