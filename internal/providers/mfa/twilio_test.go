package mfa_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/mfa"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTwilio() *mfa.Twilio {
	return mfa.NewTwilio(mfa.TwilioConfig{
		AccountSID: "test-sid",
		AuthToken:  "test-token",
		ServiceSID: "test-service",
	})
}

func TestTwilio_Name(t *testing.T) {
	assert.Equal(t, "twilio", newTwilio().Name())
}

func TestTwilio_ImplementsProvider(t *testing.T) {
	var _ mfa.Provider = newTwilio()
}

func TestTwilio_Enroll(t *testing.T) {
	tw := newTwilio()
	result, err := tw.Enroll(context.Background(), mfa.EnrollRequest{
		UserID: "user-1",
		Email:  "user@example.com",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ProviderID)
	assert.NotEmpty(t, result.QRCodeURL)
	assert.Contains(t, result.QRCodeURL, "user@example.com")
}

func TestTwilio_ConfirmEnrollment(t *testing.T) {
	tw := newTwilio()
	err := tw.ConfirmEnrollment(context.Background(), "user-1", "123456")
	assert.NoError(t, err)
}

func TestTwilio_SendChallenge(t *testing.T) {
	tw := newTwilio()
	result, err := tw.SendChallenge(context.Background(), mfa.ChallengeRequest{
		UserID:  "user-1",
		Channel: "totp",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ChallengeID)
	assert.Greater(t, result.ExpiresInS, 0)
}

func TestTwilio_VerifyCode(t *testing.T) {
	tw := newTwilio()
	result, err := tw.VerifyCode(context.Background(), mfa.VerifyRequest{
		UserID: "user-1",
		Code:   "123456",
	})
	require.NoError(t, err)
	assert.True(t, result.Valid)
}

func TestTwilio_Unenroll(t *testing.T) {
	tw := newTwilio()
	err := tw.Unenroll(context.Background(), "user-1")
	assert.NoError(t, err)
}
