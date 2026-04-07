package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupBeneficiaryService(t *testing.T) (*services.BeneficiaryService, models.Organization) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	svc := services.NewBeneficiaryService(db, reg)
	org := testutil.SeedOrganization(t, db)
	return svc, org
}

func TestBeneficiaryService_Create(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	b, err := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID,
		Name:           "Acme Corp",
		Type:           "bank",
		Currency:       "USD",
		AccountNumber:  "1234567890",
		BankName:       "Chase",
		BankCountry:    "US",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, b.ID)
	assert.Equal(t, "Acme Corp", b.Name)
	assert.Equal(t, models.BeneficiaryTypeBank, b.Type)
	assert.Equal(t, "USD", b.Currency)
	assert.Equal(t, "1234567890", *b.AccountNumber)
	assert.True(t, b.IsActive)
}

func TestBeneficiaryService_Create_Crypto(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	b, err := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID,
		Name:           "Crypto Wallet",
		Type:           "crypto",
		Currency:       "ETH",
		Address:        "0xABCDEF1234567890",
		Network:        "ethereum",
	})
	require.NoError(t, err)
	assert.Equal(t, models.BeneficiaryTypeCrypto, b.Type)
	assert.Equal(t, "0xABCDEF1234567890", *b.Address)
	assert.Equal(t, "ethereum", *b.Network)
}

func TestBeneficiaryService_Create_InvalidType(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	_, err := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID,
		Name:           "Bad Type",
		Type:           "invalid",
		Currency:       "USD",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type must be")
}

func TestBeneficiaryService_List(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID, Name: "Ben A", Type: "bank", Currency: "USD",
	})
	svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID, Name: "Ben B", Type: "bank", Currency: "EUR",
	})

	list, err := svc.List(org.ID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestBeneficiaryService_Get(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	created, _ := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID, Name: "Get Me", Type: "bank", Currency: "USD",
	})

	b, err := svc.Get(org.ID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Get Me", b.Name)
}

func TestBeneficiaryService_Get_NotFound(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	_, err := svc.Get(org.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestBeneficiaryService_Get_WrongOrg(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	created, _ := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID, Name: "Org Scoped", Type: "bank", Currency: "USD",
	})

	_, err := svc.Get("other-org", created.ID)
	assert.Error(t, err)
}

func TestBeneficiaryService_Update(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	created, _ := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID, Name: "Old Name", Type: "bank", Currency: "USD",
	})

	newName := "New Name"
	newAccount := "9999999999"
	updated, err := svc.Update(org.ID, created.ID, services.UpdateBeneficiaryInput{
		Name:          &newName,
		AccountNumber: &newAccount,
	})
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "9999999999", *updated.AccountNumber)
}

func TestBeneficiaryService_Update_NotFound(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	name := "test"
	_, err := svc.Update(org.ID, "nonexistent", services.UpdateBeneficiaryInput{Name: &name})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestBeneficiaryService_Delete(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	created, _ := svc.Create(services.CreateBeneficiaryInput{
		OrganizationID: org.ID, Name: "Delete Me", Type: "bank", Currency: "USD",
	})

	err := svc.Delete(org.ID, created.ID)
	require.NoError(t, err)

	// Should not appear in active list (soft-deleted)
	list, _ := svc.List(org.ID)
	assert.Len(t, list, 0, "should not appear in active list")
}

func TestBeneficiaryService_Delete_NotFound(t *testing.T) {
	svc, org := setupBeneficiaryService(t)

	err := svc.Delete(org.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
