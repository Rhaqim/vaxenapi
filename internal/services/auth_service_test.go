package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// registerUser is a test helper that runs the full flow:
// request access -> admin approve -> register with invite token.
func registerUser(t *testing.T, svc *services.AuthService, email, password string) *services.AuthResult {
	t.Helper()

	// Step 1: Request access
	accessReq, err := svc.RequestAccess(services.RequestAccessInput{
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Company:   "Test Corp",
		Role:      "Director",
		Country:   "US",
	})
	require.NoError(t, err)

	// Step 2: Admin approves
	approved, err := svc.ApproveAccessRequest(accessReq.ID, "admin-1")
	require.NoError(t, err)
	require.NotNil(t, approved.InviteToken)

	// Step 3: Register with invite token
	result, err := svc.Register(services.RegisterInput{
		InviteToken: *approved.InviteToken,
		FirstName:   "Test",
		LastName:    "User",
		Email:       email,
		Password:    password,
	})
	require.NoError(t, err)
	return result
}

// --- Request Access Tests ---

func TestAuthService_RequestAccess(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	req, err := svc.RequestAccess(services.RequestAccessInput{
		FirstName:    "Alice",
		LastName:     "Director",
		Email:        "alice@corp.com",
		Company:      "Alice Corp",
		Role:         "CEO",
		Country:      "US",
		Markets:      []string{"US", "EU"},
		AnnualVolume: "$1M-$10M",
		UseCase:      "Cross-border payments",
		Website:      "https://alicecorp.com",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, req.ID)
	assert.Equal(t, "alice@corp.com", req.Email)
	assert.Equal(t, models.AccessRequestPending, req.Status)
}

func TestAuthService_RequestAccess_DuplicatePending(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	input := services.RequestAccessInput{
		FirstName: "Bob", LastName: "Dir", Email: "bob@corp.com",
		Company: "Bob Corp", Role: "CFO",
	}

	_, err := svc.RequestAccess(input)
	require.NoError(t, err)

	_, err = svc.RequestAccess(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already pending")
}

func TestAuthService_GetAccessRequests(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	svc.RequestAccess(services.RequestAccessInput{
		FirstName: "A", LastName: "B", Email: "a@test.com", Company: "A Corp", Role: "CEO",
	})
	svc.RequestAccess(services.RequestAccessInput{
		FirstName: "C", LastName: "D", Email: "c@test.com", Company: "C Corp", Role: "CTO",
	})

	// All
	all, err := svc.GetAccessRequests("")
	require.NoError(t, err)
	assert.Len(t, all, 2)

	// By status
	pending, err := svc.GetAccessRequests("pending")
	require.NoError(t, err)
	assert.Len(t, pending, 2)
}

// --- Approve / Reject Tests ---

func TestAuthService_ApproveAccessRequest(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	req, _ := svc.RequestAccess(services.RequestAccessInput{
		FirstName: "App", LastName: "Rove", Email: "approve@test.com",
		Company: "Approve Corp", Role: "CEO",
	})

	approved, err := svc.ApproveAccessRequest(req.ID, "admin-1")
	require.NoError(t, err)
	assert.Equal(t, models.AccessRequestApproved, approved.Status)
	assert.NotNil(t, approved.InviteToken)
}

func TestAuthService_RejectAccessRequest(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	req, _ := svc.RequestAccess(services.RequestAccessInput{
		FirstName: "Rej", LastName: "Ect", Email: "reject@test.com",
		Company: "Reject Corp", Role: "CEO",
	})

	err := svc.RejectAccessRequest(req.ID, "admin-1", "incomplete application")
	require.NoError(t, err)

	var updated models.AccessRequest
	db.First(&updated, "id = ?", req.ID)
	assert.Equal(t, models.AccessRequestRejected, updated.Status)
	assert.Equal(t, "incomplete application", updated.ReviewNote)
}

func TestAuthService_ApproveNonPending(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	req, _ := svc.RequestAccess(services.RequestAccessInput{
		FirstName: "X", LastName: "Y", Email: "x@test.com",
		Company: "X Corp", Role: "CEO",
	})
	svc.ApproveAccessRequest(req.ID, "admin-1")

	// Try approving again
	_, err := svc.ApproveAccessRequest(req.ID, "admin-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not pending")
}

// --- Register (with invite token) Tests ---

func TestAuthService_Register(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	result := registerUser(t, svc, "register@corp.com", "strongpassword123")

	assert.NotEmpty(t, result.Token)
	assert.Equal(t, "register@corp.com", result.User.Email)
	assert.Equal(t, models.UserRoleOwner, result.User.Role)
	assert.True(t, result.User.IsDirector)
	assert.True(t, result.RequiresMFA)

	// Verify org was created
	var org models.Organization
	err := db.First(&org, "id = ?", result.User.OrganizationID).Error
	require.NoError(t, err)
	assert.Equal(t, models.KybStatusPending, org.KYBStatus)

	// Verify approval policies were created
	var policies []models.ApprovalPolicy
	db.Where("organization_id = ?", org.ID).Find(&policies)
	assert.GreaterOrEqual(t, len(policies), 3)
}

func TestAuthService_Register_InvalidInviteToken(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	_, err := svc.Register(services.RegisterInput{
		InviteToken: "bogus-token",
		FirstName:   "Bad",
		LastName:    "Token",
		Email:       "bad@test.com",
		Password:    "password1234",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired invite token")
}

func TestAuthService_Register_EmailMismatch(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	req, _ := svc.RequestAccess(services.RequestAccessInput{
		FirstName: "Mis", LastName: "Match", Email: "match@test.com",
		Company: "Match Corp", Role: "CEO",
	})
	approved, _ := svc.ApproveAccessRequest(req.ID, "admin-1")

	_, err := svc.Register(services.RegisterInput{
		InviteToken: *approved.InviteToken,
		FirstName:   "Wrong",
		LastName:    "Email",
		Email:       "wrong-email@test.com",
		Password:    "password1234",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
}

func TestAuthService_Register_TokenConsumed(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	req, _ := svc.RequestAccess(services.RequestAccessInput{
		FirstName: "Once", LastName: "Only", Email: "once@test.com",
		Company: "Once Corp", Role: "CEO",
	})
	approved, _ := svc.ApproveAccessRequest(req.ID, "admin-1")
	inviteToken := *approved.InviteToken

	// First registration succeeds
	_, err := svc.Register(services.RegisterInput{
		InviteToken: inviteToken,
		FirstName:   "Once",
		LastName:    "Only",
		Email:       "once@test.com",
		Password:    "password1234",
	})
	require.NoError(t, err)

	// Second attempt with same token fails
	_, err = svc.Register(services.RegisterInput{
		InviteToken: inviteToken,
		FirstName:   "Once",
		LastName:    "Only",
		Email:       "once@test.com",
		Password:    "password1234",
	})
	assert.Error(t, err)
}

// --- Login Tests ---

func TestAuthService_Login(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

	registerUser(t, svc, "login@corp.com", "mypassword123")

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
	svc := services.NewAuthService(db, cfg, nil)

	registerUser(t, svc, "wrong@corp.com", "correctpassword")

	_, err := svc.Login(services.LoginInput{
		Email:    "wrong@corp.com",
		Password: "wrongpassword1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_NonExistentUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	cfg := testutil.TestConfig()
	svc := services.NewAuthService(db, cfg, nil)

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
	svc := services.NewAuthService(db, cfg, nil)

	regResult := registerUser(t, svc, "mfa@corp.com", "mfapassword123")

	db.Model(&models.User{}).Where("id = ?", regResult.User.ID).Update("mfa_enabled", true)

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
	svc := services.NewAuthService(db, cfg, nil)

	regResult := registerUser(t, svc, "complete@corp.com", "password1234")

	result, err := svc.CompleteLogin(regResult.User.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, regResult.User.ID, result.User.ID)
}
