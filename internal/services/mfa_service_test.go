package services_test

import (
	"context"
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/mfa"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMFAService(t *testing.T) (*services.MFAService, *models.User) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry()
	reg.SetMFA(mfa.NewTwilio(mfa.TwilioConfig{}))
	svc := services.NewMFAService(db, reg)

	org := testutil.SeedOrganization(t, db)
	user := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	return svc, &user
}

func TestMFAService_Enroll(t *testing.T) {
	svc, user := setupMFAService(t)

	result, err := svc.EnrollMFA(context.Background(), user.ID, user.Email)
	require.NoError(t, err)
	assert.NotEmpty(t, result.QRCodeURL)
	assert.NotEmpty(t, result.ProviderID)
}

func TestMFAService_Enroll_AlreadyConfirmed(t *testing.T) {
	svc, user := setupMFAService(t)

	_, err := svc.EnrollMFA(context.Background(), user.ID, user.Email)
	require.NoError(t, err)

	// Confirm enrollment
	err = svc.ConfirmMFA(context.Background(), user.ID, "123456")
	require.NoError(t, err)

	// Try to enroll again
	_, err = svc.EnrollMFA(context.Background(), user.ID, user.Email)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already enrolled")
}

func TestMFAService_Confirm(t *testing.T) {
	svc, user := setupMFAService(t)

	_, err := svc.EnrollMFA(context.Background(), user.ID, user.Email)
	require.NoError(t, err)

	err = svc.ConfirmMFA(context.Background(), user.ID, "123456")
	require.NoError(t, err)
}

func TestMFAService_Confirm_NoEnrollment(t *testing.T) {
	svc, user := setupMFAService(t)

	err := svc.ConfirmMFA(context.Background(), user.ID, "123456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no MFA enrollment")
}

func TestMFAService_SendChallenge(t *testing.T) {
	svc, user := setupMFAService(t)

	result, err := svc.SendChallenge(context.Background(), user.ID, "totp")
	require.NoError(t, err)
	assert.NotEmpty(t, result.ChallengeID)
	assert.Greater(t, result.ExpiresInS, 0)
}

func TestMFAService_VerifyCode(t *testing.T) {
	svc, user := setupMFAService(t)

	err := svc.VerifyCode(context.Background(), user.ID, "123456")
	assert.NoError(t, err)
}

func TestMFAService_Unenroll(t *testing.T) {
	svc, user := setupMFAService(t)

	// Enroll and confirm
	_, err := svc.EnrollMFA(context.Background(), user.ID, user.Email)
	require.NoError(t, err)
	err = svc.ConfirmMFA(context.Background(), user.ID, "123456")
	require.NoError(t, err)

	// Unenroll
	err = svc.UnenrollMFA(context.Background(), user.ID)
	require.NoError(t, err)

	// Can enroll again
	_, err = svc.EnrollMFA(context.Background(), user.ID, user.Email)
	assert.NoError(t, err)
}

func TestMFAService_NilProvider(t *testing.T) {
	db := testutil.SetupTestDB(t)
	reg := providers.NewRegistry() // no MFA provider set
	svc := services.NewMFAService(db, reg)

	_, err := svc.EnrollMFA(context.Background(), "user-1", "test@test.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}
