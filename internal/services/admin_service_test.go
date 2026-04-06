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

func TestAdminService_Settings(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)

	// Create setting
	err := svc.UpsertSetting("max_payout", "50000", "limits", "admin-1")
	require.NoError(t, err)

	err = svc.UpsertSetting("min_payout", "10", "limits", "admin-1")
	require.NoError(t, err)

	// Get all
	settings, err := svc.GetSettings("")
	require.NoError(t, err)
	assert.Len(t, settings, 2)

	// Get by category
	settings, err = svc.GetSettings("limits")
	require.NoError(t, err)
	assert.Len(t, settings, 2)

	// Update existing
	err = svc.UpsertSetting("max_payout", "100000", "limits", "admin-2")
	require.NoError(t, err)

	settings, err = svc.GetSettings("limits")
	require.NoError(t, err)
	for _, s := range settings {
		if s.Key == "max_payout" {
			assert.Equal(t, "100000", s.Value)
			assert.Equal(t, "admin-2", s.UpdatedBy)
		}
	}
}

func TestAdminService_FeatureFlags(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)

	// Create flag (disabled)
	err := svc.SetFeatureFlag("crypto_trading", false, "Enable crypto trading", "admin-1")
	require.NoError(t, err)

	assert.False(t, svc.IsFeatureEnabled("crypto_trading"))

	// Enable it
	err = svc.SetFeatureFlag("crypto_trading", true, "Enable crypto trading", "admin-1")
	require.NoError(t, err)

	assert.True(t, svc.IsFeatureEnabled("crypto_trading"))

	// List all
	flags, err := svc.GetFeatureFlags()
	require.NoError(t, err)
	assert.Len(t, flags, 1)
	assert.Equal(t, "crypto_trading", flags[0].Name)
}

func TestAdminService_ExchangeRates(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)

	rate := decimal.NewFromFloat(0.85)
	spread := decimal.NewFromFloat(0.005)

	// Create rate
	err := svc.UpsertExchangeRate("USD", "EUR", rate, spread, "admin-1")
	require.NoError(t, err)

	rates, err := svc.GetExchangeRates()
	require.NoError(t, err)
	assert.Len(t, rates, 1)
	assert.True(t, rates[0].Rate.Equal(rate))
	assert.True(t, rates[0].Spread.Equal(spread))

	// Update rate
	newRate := decimal.NewFromFloat(0.90)
	err = svc.UpsertExchangeRate("USD", "EUR", newRate, spread, "admin-2")
	require.NoError(t, err)

	rates, err = svc.GetExchangeRates()
	require.NoError(t, err)
	assert.Len(t, rates, 1)
	assert.True(t, rates[0].Rate.Equal(newRate))
}

func TestAdminService_ApproveOrganization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)
	org := testutil.SeedOrganization(t, db)

	err := svc.ApproveOrganization(org.ID)
	require.NoError(t, err)

	var updated models.Organization
	db.First(&updated, "id = ?", org.ID)
	assert.Equal(t, models.KybStatusApproved, updated.KYBStatus)
}

func TestAdminService_RejectOrganization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)
	org := testutil.SeedOrganization(t, db)

	err := svc.RejectOrganization(org.ID, "incomplete documents")
	require.NoError(t, err)

	var updated models.Organization
	db.First(&updated, "id = ?", org.ID)
	assert.Equal(t, models.KybStatusRejected, updated.KYBStatus)
}

func TestAdminService_ApproveNonExistent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)

	err := svc.ApproveOrganization("nonexistent-id")
	assert.Error(t, err)
}

func TestAdminService_GetAllUsers(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)
	org := testutil.SeedOrganization(t, db)
	testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	users, total, err := svc.GetAllUsers(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, users, 1)
}

func TestAdminService_GetAllOrganizations(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewAdminService(db)
	testutil.SeedOrganization(t, db)

	orgs, total, err := svc.GetAllOrganizations(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, orgs, 1)
}
