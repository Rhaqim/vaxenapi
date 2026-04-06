package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/payment"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gorm.io/gorm"
)

func setupPaymentService(t *testing.T) (*services.PaymentService, *gorm.DB, models.Organization, models.User) {
	t.Helper()
	db := testutil.SetupTestDB(t)

	// Also migrate Beneficiary for payment tests
	db.AutoMigrate(&models.Beneficiary{}, &models.Payout{})

	reg := providers.NewRegistry()
	reg.SetPayment(payment.NewCircle(payment.CircleConfig{}))
	approval := services.NewApprovalService(db)
	svc := services.NewPaymentService(db, reg, approval)
	org := testutil.SeedOrganization(t, db)
	user := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)
	return svc, db, org, user
}

func seedBeneficiary(t *testing.T, db *gorm.DB, orgID string) models.Beneficiary {
	t.Helper()
	b := models.Beneficiary{
		OrganizationID: orgID,
		Name:           "Test Beneficiary",
	}
	if err := db.Create(&b).Error; err != nil {
		t.Fatalf("failed to seed beneficiary: %v", err)
	}
	return b
}

func TestPaymentService_InitiatePayment_SingleApproval(t *testing.T) {
	svc, db, org, user := setupPaymentService(t)
	beneficiary := seedBeneficiary(t, db, org.ID)

	// Default approval policy is 1 — should execute immediately
	payout, err := svc.InitiatePayment(t.Context(), services.SendPaymentInput{
		OrganizationID: org.ID,
		InitiatedByID:  user.ID,
		BeneficiaryID:  beneficiary.ID,
		Amount:         decimal.NewFromFloat(500),
		Currency:       "USD",
		Reference:      "PAY-001",
		Description:    "Test payment",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, payout.ID)
	assert.Equal(t, models.PayoutStatusProcessing, payout.Status)
}

func TestPaymentService_InitiatePayment_MultiApproval(t *testing.T) {
	svc, db, org, user := setupPaymentService(t)
	beneficiary := seedBeneficiary(t, db, org.ID)

	// Set policy to require 2 approvals
	approval := services.NewApprovalService(db)
	approval.UpdatePolicy(org.ID, models.ApprovalActionPayout, 2)

	payout, err := svc.InitiatePayment(t.Context(), services.SendPaymentInput{
		OrganizationID: org.ID,
		InitiatedByID:  user.ID,
		BeneficiaryID:  beneficiary.ID,
		Amount:         decimal.NewFromFloat(1000),
		Currency:       "USD",
		Reference:      "PAY-002",
	})
	require.NoError(t, err)
	assert.Equal(t, models.PayoutStatusPending, payout.Status, "should stay pending until approved")

	// Verify approval request was created
	requests, err := approval.GetPendingRequests(org.ID)
	require.NoError(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, payout.ID, requests[0].ResourceID)
}

func TestPaymentService_InitiatePayment_BeneficiaryNotFound(t *testing.T) {
	svc, _, org, user := setupPaymentService(t)

	_, err := svc.InitiatePayment(t.Context(), services.SendPaymentInput{
		OrganizationID: org.ID,
		InitiatedByID:  user.ID,
		BeneficiaryID:  "nonexistent",
		Amount:         decimal.NewFromFloat(100),
		Currency:       "USD",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestPaymentService_GetPayouts(t *testing.T) {
	svc, db, org, user := setupPaymentService(t)
	beneficiary := seedBeneficiary(t, db, org.ID)

	svc.InitiatePayment(t.Context(), services.SendPaymentInput{
		OrganizationID: org.ID,
		InitiatedByID:  user.ID,
		BeneficiaryID:  beneficiary.ID,
		Amount:         decimal.NewFromFloat(100),
		Currency:       "USD",
		Reference:      "P1",
	})

	payouts, err := svc.GetPayouts(org.ID)
	require.NoError(t, err)
	assert.Len(t, payouts, 1)
}

func TestPaymentService_GetPaymentStatus(t *testing.T) {
	svc, db, org, user := setupPaymentService(t)
	beneficiary := seedBeneficiary(t, db, org.ID)

	payout, err := svc.InitiatePayment(t.Context(), services.SendPaymentInput{
		OrganizationID: org.ID,
		InitiatedByID:  user.ID,
		BeneficiaryID:  beneficiary.ID,
		Amount:         decimal.NewFromFloat(200),
		Currency:       "USD",
		Reference:      "P2",
	})
	require.NoError(t, err)

	status, err := svc.GetPaymentStatus(t.Context(), payout.ID, org.ID)
	require.NoError(t, err)
	assert.Equal(t, payout.ID, status.ID)
}

func TestPaymentService_NilProvider(t *testing.T) {
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.Beneficiary{}, &models.Payout{})
	reg := providers.NewRegistry()
	approval := services.NewApprovalService(db)
	svc := services.NewPaymentService(db, reg, approval)

	_, err := svc.GetFees(t.Context(), decimal.NewFromInt(100), "USD", "US", "wire")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}
