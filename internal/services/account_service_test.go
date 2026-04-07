package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAccountService(t *testing.T) (*services.AccountService, models.Organization) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.AccountNumber{})
	svc := services.NewAccountService(db)
	org := testutil.SeedOrganization(t, db)
	return svc, org
}

func TestAccountService_Create(t *testing.T) {
	svc, org := setupAccountService(t)

	a, err := svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID,
		Name:           "Main USD Account",
		Currency:       "USD",
		Type:           "ach",
		AccountNumber:  "1234567890",
		BankName:       "Chase",
		BankCountry:    "US",
		RoutingNumber:  "021000021",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, a.ID)
	assert.Equal(t, "Main USD Account", a.Name)
	assert.Equal(t, models.AccountTypeACH, a.Type)
	assert.Equal(t, "1234567890", a.AccountNumber)
	assert.Equal(t, "021000021", *a.RoutingNumber)
	assert.True(t, a.IsActive)
}

func TestAccountService_Create_IBAN(t *testing.T) {
	svc, org := setupAccountService(t)

	a, err := svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID,
		Name:           "EUR IBAN",
		Currency:       "EUR",
		Type:           "iban",
		AccountNumber:  "DE89370400440532013000",
		BankName:       "Deutsche Bank",
		BankCountry:    "DE",
	})
	require.NoError(t, err)
	assert.Equal(t, models.AccountTypeIBAN, a.Type)
}

func TestAccountService_Create_InvalidType(t *testing.T) {
	svc, org := setupAccountService(t)

	_, err := svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID,
		Name:           "Bad",
		Currency:       "USD",
		Type:           "invalid",
		AccountNumber:  "123",
		BankName:       "Bank",
		BankCountry:    "US",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type must be")
}

func TestAccountService_List(t *testing.T) {
	svc, org := setupAccountService(t)

	svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID, Name: "A", Currency: "USD", Type: "ach",
		AccountNumber: "111", BankName: "Bank A", BankCountry: "US",
	})
	svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID, Name: "B", Currency: "EUR", Type: "iban",
		AccountNumber: "222", BankName: "Bank B", BankCountry: "DE",
	})

	list, err := svc.List(org.ID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestAccountService_Get(t *testing.T) {
	svc, org := setupAccountService(t)

	created, _ := svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID, Name: "Get Me", Currency: "USD", Type: "ach",
		AccountNumber: "333", BankName: "Bank", BankCountry: "US",
	})

	a, err := svc.Get(org.ID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Get Me", a.Name)
}

func TestAccountService_Get_NotFound(t *testing.T) {
	svc, org := setupAccountService(t)

	_, err := svc.Get(org.ID, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAccountService_Get_WrongOrg(t *testing.T) {
	svc, org := setupAccountService(t)

	created, _ := svc.Create(services.CreateAccountInput{
		OrganizationID: org.ID, Name: "Scoped", Currency: "USD", Type: "ach",
		AccountNumber: "444", BankName: "Bank", BankCountry: "US",
	})

	_, err := svc.Get("other-org", created.ID)
	assert.Error(t, err)
}
