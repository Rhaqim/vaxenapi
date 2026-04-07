package database

import (
	"log"

	"vaxen/api/internal/config"
	"vaxen/api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the database connection
func Connect(cfg *config.Config) error {
	var err error

	gormConfig := &gorm.Config{}

	if cfg.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return err
	}

	log.Println("✅ Database connection established")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	return DB.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Wallet{},
		&models.WalletTransaction{},
		&models.AccountNumber{},
		&models.Beneficiary{},
		&models.Deposit{},
		&models.Payout{},
		&models.PayoutBatch{},
		&models.ConversionOrder{},
		&models.LimitOrder{},
		&models.Journal{},
		&models.LedgerEntry{},
		&models.ReconciliationRun{},
		&models.StatementFile{},
		&models.ComplianceCase{},
		&models.ScreeningResult{},
		&models.AuditLog{},
		&models.WebhookEvent{},
		// New models
		&models.ApprovalPolicy{},
		&models.ApprovalRequest{},
		&models.ApprovalVote{},
		&models.PlatformSetting{},
		&models.FeatureFlag{},
		&models.ExchangeRate{},
		&models.MFAEnrollment{},
		&models.MFAChallenge{},
		&models.Web3Wallet{},
		&models.AccessRequest{},
	)
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
