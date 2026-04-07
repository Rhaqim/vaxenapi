package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportService_GetStatements_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.StatementFile{})
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	stmts, err := svc.GetStatements(org.ID)
	require.NoError(t, err)
	assert.Len(t, stmts, 0)
}

func TestReportService_GetStatement_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.StatementFile{})
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	_, err := svc.GetStatement(org.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestReportService_GetAuditLogs(t *testing.T) {
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.AuditLog{})
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	// Seed some audit logs
	db.Create(&models.AuditLog{OrganizationID: org.ID, Action: "user.login", Resource: "user"})
	db.Create(&models.AuditLog{OrganizationID: org.ID, Action: "payout.created", Resource: "payout"})

	logs, total, err := svc.GetAuditLogs(org.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, logs, 2)
}

func TestReportService_GetAuditLogs_Pagination(t *testing.T) {
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.AuditLog{})
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	for i := 0; i < 5; i++ {
		db.Create(&models.AuditLog{OrganizationID: org.ID, Action: "test", Resource: "test"})
	}

	logs, total, err := svc.GetAuditLogs(org.ID, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, logs, 2, "should return only page size")
}

func TestReportService_GetBalanceReport(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	db.Create(&models.Wallet{
		OrganizationID: org.ID, Type: models.WalletTypeFiat,
		Currency: "USD", Balance: decimal.NewFromInt(1000),
		AvailableBalance: decimal.NewFromInt(950), IsActive: true,
	})
	db.Create(&models.Wallet{
		OrganizationID: org.ID, Type: models.WalletTypeFiat,
		Currency: "EUR", Balance: decimal.NewFromInt(500),
		AvailableBalance: decimal.NewFromInt(500), IsActive: true,
	})

	entries, err := svc.GetBalanceReport(org.ID)
	require.NoError(t, err)
	assert.Len(t, entries, 2)

	currencies := map[string]bool{}
	for _, e := range entries {
		currencies[e.Currency] = true
	}
	assert.True(t, currencies["USD"])
	assert.True(t, currencies["EUR"])
}

func TestReportService_GetBalanceReport_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	entries, err := svc.GetBalanceReport(org.ID)
	require.NoError(t, err)
	assert.Len(t, entries, 0)
}

func TestReportService_GetTransactionReport(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	// No transactions yet
	report, err := svc.GetTransactionReport(org.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), report.TotalTransactions)
}

func TestReportService_GetFxPnlReport(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewReportService(db)
	org := testutil.SeedOrganization(t, db)

	// Seed a conversion order
	db.Create(&models.ConversionOrder{
		OrganizationID: org.ID,
		FromAmount:     decimal.NewFromInt(1000),
		ToAmount:       decimal.NewFromInt(850),
		Rate:           decimal.NewFromFloat(0.85),
		Fee:            decimal.NewFromFloat(1.50),
		Status:         models.OrderStatusCompleted,
	})

	report, err := svc.GetFxPnlReport(org.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), report.TotalConversions)
	assert.True(t, report.TotalFees.GreaterThan(decimal.Zero))
}
