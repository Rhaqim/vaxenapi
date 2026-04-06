package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DatabaseURL string

	// JWT / Auth
	JWTSecret              string
	JWTExpiration          time.Duration // Access token lifetime (short — e.g. 15m)
	RefreshTokenExpiration time.Duration // Refresh token lifetime (e.g. 7d)
	SecureCookies          bool          // true in production (requires HTTPS)

	// CORS
	FrontendURL string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Logging
	LogLevel string

	// Provider configuration
	Providers ProviderConfig
}

// ProviderConfig holds configuration for all third-party providers.
type ProviderConfig struct {
	// KYC/KYB
	KYCProvider     string // "sumsub", "onfido"
	SumSubAPIKey    string
	SumSubAPISecret string
	SumSubBaseURL   string

	// MFA
	MFAProvider      string // "twilio", "totp"
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioServiceSID string

	// Wallet / Key Management
	WalletProvider  string // "aws_kms", "fireblocks"
	AWSRegion       string
	AWSKMSKeyPolicy string

	// Exchange
	ExchangeProvider string // "internal", "binance", "kraken"

	// Payment
	PaymentProvider string // "circle", "wise", "stripe"
	CircleAPIKey    string
	CircleBaseURL   string
}

// Load returns the application configuration loaded from environment variables
func Load() *Config {

	return &Config{
		Port:                   getEnv("PORT", "8080").String(),
		Environment:            getEnv("ENVIRONMENT", "development").String(),
		DatabaseURL:            getEnv("DATABASE_URL", "").String(),
		JWTSecret:              getEnv("JWT_SECRET", "change-this-secret").String(),
		JWTExpiration:          getEnv("JWT_EXPIRATION", "15m").AsDuration(),
		RefreshTokenExpiration: getEnv("REFRESH_TOKEN_EXPIRATION", "168h").AsDuration(), // 7 days
		SecureCookies:          getEnv("ENVIRONMENT", "development").String() == "production",
		FrontendURL:            getEnv("FRONTEND_URL", "http://localhost:3000").String(),
		RedisHost:              getEnv("REDIS_HOST", "localhost").String(),
		RedisPort:              getEnv("REDIS_PORT", "6379").String(),
		RedisPassword:          getEnv("REDIS_PASSWORD", "").String(),
		LogLevel:               getEnv("LOG_LEVEL", "info").String(),

		Providers: ProviderConfig{
			// KYC
			KYCProvider:     getEnv("KYC_PROVIDER", "sumsub").String(),
			SumSubAPIKey:    getEnv("SUMSUB_API_KEY", "").String(),
			SumSubAPISecret: getEnv("SUMSUB_API_SECRET", "").String(),
			SumSubBaseURL:   getEnv("SUMSUB_BASE_URL", "https://api.sumsub.com").String(),

			// MFA
			MFAProvider:      getEnv("MFA_PROVIDER", "twilio").String(),
			TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", "").String(),
			TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", "").String(),
			TwilioServiceSID: getEnv("TWILIO_SERVICE_SID", "").String(),

			// Wallet
			WalletProvider:  getEnv("WALLET_PROVIDER", "aws_kms").String(),
			AWSRegion:       getEnv("AWS_REGION", "us-east-1").String(),
			AWSKMSKeyPolicy: getEnv("AWS_KMS_KEY_POLICY", "").String(),

			// Exchange
			ExchangeProvider: getEnv("EXCHANGE_PROVIDER", "internal").String(),

			// Payment
			PaymentProvider: getEnv("PAYMENT_PROVIDER", "circle").String(),
			CircleAPIKey:    getEnv("CIRCLE_API_KEY", "").String(),
			CircleBaseURL:   getEnv("CIRCLE_BASE_URL", "https://api.circle.com").String(),
		},
	}
}

func getEnv(key, defaultValue string) *EnvValues {
	if value := os.Getenv(key); value != "" {
		return &EnvValues{value: value}
	}
	return &EnvValues{value: defaultValue}
}

type EnvValues struct {
	value string
}

func (e *EnvValues) String() string {
	return e.value
}

func (e *EnvValues) AsInt() (int, error) {
	// Parse as int, return error if unable
	var result int
	_, err := fmt.Sscanf(e.value, "%d", &result)
	return result, err
}

func (e *EnvValues) AsBool() (bool, error) {
	// Parse as bool, return error if unable
	return strconv.ParseBool(e.value)
}

func (e *EnvValues) AsDuration() time.Duration {
	// Parse as duration, return error if unable
	val, _ := time.ParseDuration(e.value)
	return val
}
