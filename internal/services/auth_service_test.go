package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	result, err := svc.Register(services.RegisterInput{
		Email:              "alice@corp.com",
		Password:           "strongpassword123",
		FirstName:          "Alice",
		LastName:           "Director",
		CompanyName:        "Alice Corp",
		LegalName:          "Alice Corporation Ltd",
		RegistrationNumber: "AC-001",
		Country:            "US",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "alice@corp.com", result.User.Email)
	assert.Equal(t, models.UserRoleOwner, result.User.Role)
	assert.True(t, result.User.IsDirector)
	assert.True(t, result.RequiresMFA, "directors must set up MFA")

	// Verify org was created
	var org models.Organization
	err = db.First(&org, "id = ?", result.User.OrganizationID).Error
	require.NoError(t, err)
	assert.Equal(t, "Alice Corporation Ltd", org.LegalName)
	assert.Equal(t, models.KybStatusPending, org.KYBStatus)

	// Verify default approval policies were created
	var policies []models.ApprovalPolicy
	db.Where("organization_id = ?", org.ID).Find(&policies)
	assert.GreaterOrEqual(t, len(policies), 3)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	input := services.RegisterInput{
		Email:              "bob@corp.com",
		Password:           "strongpassword123",
		FirstName:          "Bob",
		LastName:           "Director",
		CompanyName:        "Bob Corp",
		LegalName:          "Bob Corp Ltd",
		RegistrationNumber: "BC-001",
		Country:            "GB",
	}

	_, err := svc.Register(input)
	require.NoError(t, err)

	_, err = svc.Register(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestAuthService_Login(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	// Register first
	_, err := svc.Register(services.RegisterInput{
		Email:              "login@corp.com",
		Password:           "mypassword123",
		FirstName:          "Login",
		LastName:           "User",
		CompanyName:        "Login Corp",
		LegalName:          "Login Corp Ltd",
		RegistrationNumber: "LC-001",
		Country:            "US",
	})
	require.NoError(t, err)

	// Login with correct credentials (MFA not enabled yet)
	result, err := svc.Login(services.LoginInput{
		Email:    "login@corp.com",
		Password: "mypassword123",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "login@corp.com", result.User.Email)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	_, err := svc.Register(services.RegisterInput{
		Email:              "wrong@corp.com",
		Password:           "correctpassword",
		FirstName:          "Wrong",
		LastName:           "Pass",
		CompanyName:        "WP Corp",
		LegalName:          "WP Corp Ltd",
		RegistrationNumber: "WP-001",
		Country:            "US",
	})
	require.NoError(t, err)

	_, err = svc.Login(services.LoginInput{
		Email:    "wrong@corp.com",
		Password: "wrongpassword1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_NonExistentUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	_, err := svc.Login(services.LoginInput{
		Email:    "nobody@corp.com",
		Password: "whatever12345",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_MFAEnabled(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	// Register user
	regResult, err := svc.Register(services.RegisterInput{
		Email:              "mfa@corp.com",
		Password:           "mfapassword123",
		FirstName:          "MFA",
		LastName:           "User",
		CompanyName:        "MFA Corp",
		LegalName:          "MFA Corp Ltd",
		RegistrationNumber: "MFA-001",
		Country:            "US",
	})
	require.NoError(t, err)

	// Enable MFA on the user
	db.Model(&models.User{}).Where("id = ?", regResult.User.ID).Update("mfa_enabled", true)

	// Login should require MFA
	result, err := svc.Login(services.LoginInput{
		Email:    "mfa@corp.com",
		Password: "mfapassword123",
	})
	require.NoError(t, err)
	assert.True(t, result.RequiresMFA)
	assert.Empty(t, result.Token, "should not issue token before MFA")
}

func TestAuthService_CompleteLogin(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg)

	regResult, err := svc.Register(services.RegisterInput{
		Email:              "complete@corp.com",
		Password:           "password1234",
		FirstName:          "Complete",
		LastName:           "Login",
		CompanyName:        "CL Corp",
		LegalName:          "CL Corp Ltd",
		RegistrationNumber: "CL-001",
		Country:            "US",
	})
	require.NoError(t, err)

	result, err := svc.CompleteLogin(regResult.User.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, regResult.User.ID, result.User.ID)
}
