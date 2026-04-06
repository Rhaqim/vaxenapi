package kyc_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/kyc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSumSub() *kyc.SumSub {
	return kyc.NewSumSub(kyc.SumSubConfig{
		APIKey:    "test-key",
		APISecret: "test-secret",
	})
}

func TestSumSub_Name(t *testing.T) {
	assert.Equal(t, "sumsub", newSumSub().Name())
}

func TestSumSub_ImplementsProvider(t *testing.T) {
	var _ kyc.Provider = newSumSub()
}

func TestSumSub_SubmitKYB(t *testing.T) {
	s := newSumSub()
	result, err := s.SubmitKYB(context.Background(), kyc.KYBSubmission{
		OrganizationID:     "org-123",
		LegalName:          "Test Corp",
		RegistrationNumber: "REG-001",
		Country:            "US",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ReferenceID)
	assert.Equal(t, kyc.StatusPending, result.Status)
}

func TestSumSub_GetKYBStatus(t *testing.T) {
	s := newSumSub()
	result, err := s.GetKYBStatus(context.Background(), "ref-123")
	require.NoError(t, err)
	assert.Equal(t, "ref-123", result.ReferenceID)
	assert.Equal(t, kyc.StatusPending, result.Status)
}

func TestSumSub_SubmitKYC(t *testing.T) {
	s := newSumSub()
	result, err := s.SubmitKYC(context.Background(), kyc.KYCSubmission{
		UserID:    "user-123",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.ReferenceID)
	assert.Equal(t, kyc.StatusPending, result.Status)
}

func TestSumSub_GetKYCStatus(t *testing.T) {
	s := newSumSub()
	result, err := s.GetKYCStatus(context.Background(), "ref-456")
	require.NoError(t, err)
	assert.Equal(t, "ref-456", result.ReferenceID)
}

func TestSumSub_GenerateVerificationURL(t *testing.T) {
	s := newSumSub()
	url, err := s.GenerateVerificationURL(context.Background(), "applicant-1", "https://example.com/done")
	require.NoError(t, err)
	assert.Contains(t, url, "applicant-1")
	assert.Contains(t, url, "example.com/done")
}
