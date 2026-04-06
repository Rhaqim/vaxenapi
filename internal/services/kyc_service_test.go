package services_test

import (
	"context"
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupKYCService(t *testing.T) (*services.KYCService, models.Organization) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	svc := services.NewKYCService(db, reg)
	org := testutil.SeedOrganization(t, db)
	return svc, org
}

func TestKYCService_SubmitKYB(t *testing.T) {
	svc, org := setupKYCService(t)

	result, err := svc.SubmitKYB(context.Background(), services.KYBSubmitInput{
		OrganizationID:     org.ID,
		LegalName:          "Test Corp Ltd",
		RegistrationNumber: "REG-123",
		Country:            "US",
		Address: kyc.Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "US",
		},
		Directors: []kyc.DirectorInfo{
			{FirstName: "John", LastName: "Doe", Email: "john@test.com"},
		},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ReferenceID)
	assert.Equal(t, kyc.StatusPending, result.Status)
}

func TestKYCService_SubmitKYB_AlreadyApproved(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	svc := services.NewKYCService(db, reg)

	org := testutil.SeedOrganization(t, db)
	db.Model(&org).Update("kyb_status", models.KybStatusApproved)

	_, err := svc.SubmitKYB(context.Background(), services.KYBSubmitInput{
		OrganizationID: org.ID,
		LegalName:      "Test",
		Country:        "US",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already approved")
}

func TestKYCService_GetKYBStatus_NoReference(t *testing.T) {
	svc, org := setupKYCService(t)

	result, err := svc.GetKYBStatus(context.Background(), org.ID)
	require.NoError(t, err)
	assert.Equal(t, kyc.VerificationStatus("pending"), result.Status)
}

func TestKYCService_GetKYBStatus_WithReference(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	reg.SetKYC(kyc.NewSumSub(kyc.SumSubConfig{}))
	svc := services.NewKYCService(db, reg)

	org := testutil.SeedOrganization(t, db)
	db.Model(&org).Update("kyb_reference_id", "ref-xyz")

	result, err := svc.GetKYBStatus(context.Background(), org.ID)
	require.NoError(t, err)
	assert.Equal(t, "ref-xyz", result.ReferenceID)
}

func TestKYCService_SubmitKYC(t *testing.T) {
	svc, _ := setupKYCService(t)

	result, err := svc.SubmitKYC(context.Background(), kyc.KYCSubmission{
		UserID:    "user-1",
		FirstName: "Jane",
		LastName:  "Director",
		Email:     "jane@test.com",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ReferenceID)
}

func TestKYCService_NilProvider(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	svc := services.NewKYCService(db, reg)

	_, err := svc.SubmitKYB(context.Background(), services.KYBSubmitInput{OrganizationID: "org-1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}

func TestKYCService_NonExistentOrg(t *testing.T) {
	svc, _ := setupKYCService(t)

	_, err := svc.SubmitKYB(context.Background(), services.KYBSubmitInput{
		OrganizationID: "nonexistent",
		LegalName:      "Ghost Corp",
		Country:        "US",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
