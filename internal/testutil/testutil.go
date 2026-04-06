package testutil

import (
	"testing"

	"vaxen/api/internal/config"
	"vaxen/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database with all migrations applied.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Wallet{},
		&models.WalletTransaction{},
		&models.ApprovalPolicy{},
		&models.ApprovalRequest{},
		&models.ApprovalVote{},
		&models.PlatformSetting{},
		&models.FeatureFlag{},
		&models.ExchangeRate{},
		&models.MFAEnrollment{},
		&models.MFAChallenge{},
		&models.Web3Wallet{},
		&models.Beneficiary{},
		&models.Payout{},
		&models.ConversionOrder{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return db
}

// TestConfig returns a config suitable for testing.
func TestConfig() *config.Config {
	return &config.Config{
		Port:          "8080",
		Environment:   "test",
		JWTSecret:     "test-secret-key-for-testing-only",
		JWTExpiration: 24 * 60 * 60 * 1000000000, // 24h in nanoseconds as Duration
		Providers: config.ProviderConfig{
			KYCProvider:      "sumsub",
			MFAProvider:      "twilio",
			WalletProvider:   "aws_kms",
			ExchangeProvider: "internal",
			PaymentProvider:  "circle",
		},
	}
}

// SeedOrganization creates a test organization.
func SeedOrganization(t *testing.T, db *gorm.DB) models.Organization {
	t.Helper()
	org := models.Organization{
		Name:               "Test Corp",
		LegalName:          "Test Corporation Ltd",
		RegistrationNumber: "TC-12345",
		Country:            "US",
		KYBStatus:          models.KybStatusPending,
	}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("failed to seed org: %v", err)
	}
	return org
}

// SeedUser creates a test user belonging to the given org.
func SeedUser(t *testing.T, db *gorm.DB, orgID string, role models.UserRole, isDirector bool) models.User {
	t.Helper()
	user := models.User{
		Email:          "testuser-" + uuid.New().String()[:8] + "@example.com",
		FirstName:      "Test",
		LastName:       "User",
		PasswordHash:   "$2a$10$abcdefghijklmnopqrstuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu", // placeholder
		Role:           role,
		OrganizationID: orgID,
		IsDirector:     isDirector,
		IsActive:       true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return user
}
