package services_test

import (
	"context"
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovalService_GetRequiredApprovals_Default(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)

	// No policy exists — should default to 1
	count, err := svc.GetRequiredApprovals("nonexistent-org", models.ApprovalActionPayout)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestApprovalService_UpdatePolicy(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)
	org := testutil.SeedOrganization(t, db)

	// Create policy
	err := svc.UpdatePolicy(org.ID, models.ApprovalActionPayout, 2)
	require.NoError(t, err)

	count, err := svc.GetRequiredApprovals(org.ID, models.ApprovalActionPayout)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Update existing policy
	err = svc.UpdatePolicy(org.ID, models.ApprovalActionPayout, 3)
	require.NoError(t, err)

	count, err = svc.GetRequiredApprovals(org.ID, models.ApprovalActionPayout)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestApprovalService_CreateAndVote(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)
	org := testutil.SeedOrganization(t, db)
	director1 := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	// Create a second director with unique email
	director2 := models.User{
		Email:          "director2@test.com",
		FirstName:      "Dir",
		LastName:       "Two",
		PasswordHash:   "hash",
		Role:           models.UserRoleManager,
		OrganizationID: org.ID,
		IsDirector:     true,
		IsActive:       true,
	}
	db.Create(&director2)

	// Create approval request requiring 2 approvals
	req, err := svc.CreateRequest(context.Background(), services.CreateApprovalInput{
		OrganizationID: org.ID,
		ActionType:     models.ApprovalActionPayout,
		ResourceID:     "payout-123",
		ResourceType:   "payout",
		RequestedByID:  director1.ID,
		RequiredCount:  2,
		Payload:        []byte(`{"amount": "1000"}`),
	})
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalStatusPending, req.Status)
	assert.Equal(t, 2, req.RequiredCount)
	assert.Equal(t, 0, req.CurrentCount)

	// First vote
	err = svc.Vote(context.Background(), req.ID, director1.ID, models.ApprovalStatusApproved, "looks good")
	require.NoError(t, err)

	// Check it's still pending
	updated, err := svc.GetRequest(req.ID, org.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalStatusPending, updated.Status)
	assert.Equal(t, 1, updated.CurrentCount)

	// Second vote — should complete
	err = svc.Vote(context.Background(), req.ID, director2.ID, models.ApprovalStatusApproved, "approved")
	require.NoError(t, err)

	final, err := svc.GetRequest(req.ID, org.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalStatusApproved, final.Status)
	assert.Equal(t, 2, final.CurrentCount)
	assert.NotNil(t, final.CompletedAt)
}

func TestApprovalService_Vote_Rejection(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)
	org := testutil.SeedOrganization(t, db)
	director := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	req, err := svc.CreateRequest(context.Background(), services.CreateApprovalInput{
		OrganizationID: org.ID,
		ActionType:     models.ApprovalActionPayout,
		ResourceID:     "payout-456",
		ResourceType:   "payout",
		RequestedByID:  director.ID,
		RequiredCount:  2,
		Payload:        []byte(`{}`),
	})
	require.NoError(t, err)

	// Reject
	err = svc.Vote(context.Background(), req.ID, director.ID, models.ApprovalStatusRejected, "too expensive")
	require.NoError(t, err)

	updated, err := svc.GetRequest(req.ID, org.ID)
	require.NoError(t, err)
	assert.Equal(t, models.ApprovalStatusRejected, updated.Status)
}

func TestApprovalService_Vote_DuplicateVote(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)
	org := testutil.SeedOrganization(t, db)
	director := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	req, err := svc.CreateRequest(context.Background(), services.CreateApprovalInput{
		OrganizationID: org.ID,
		ActionType:     models.ApprovalActionPayout,
		ResourceID:     "payout-789",
		ResourceType:   "payout",
		RequestedByID:  director.ID,
		RequiredCount:  2,
		Payload:        []byte(`{}`),
	})
	require.NoError(t, err)

	err = svc.Vote(context.Background(), req.ID, director.ID, models.ApprovalStatusApproved, "ok")
	require.NoError(t, err)

	err = svc.Vote(context.Background(), req.ID, director.ID, models.ApprovalStatusApproved, "ok again")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already voted")
}

func TestApprovalService_Vote_InsufficientRole(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)
	org := testutil.SeedOrganization(t, db)
	director := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	viewer := models.User{
		Email:          "viewer@test.com",
		FirstName:      "View",
		LastName:       "Only",
		PasswordHash:   "hash",
		Role:           models.UserRoleViewer,
		OrganizationID: org.ID,
		IsDirector:     false,
		IsActive:       true,
	}
	db.Create(&viewer)

	req, err := svc.CreateRequest(context.Background(), services.CreateApprovalInput{
		OrganizationID: org.ID,
		ActionType:     models.ApprovalActionPayout,
		ResourceID:     "payout-abc",
		ResourceType:   "payout",
		RequestedByID:  director.ID,
		RequiredCount:  1,
		Payload:        []byte(`{}`),
	})
	require.NoError(t, err)

	err = svc.Vote(context.Background(), req.ID, viewer.ID, models.ApprovalStatusApproved, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")
}

func TestApprovalService_GetPendingRequests(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewApprovalService(db)
	org := testutil.SeedOrganization(t, db)
	director := testutil.SeedUser(t, db, org.ID, models.UserRoleOwner, true)

	// Create 2 requests
	svc.CreateRequest(context.Background(), services.CreateApprovalInput{
		OrganizationID: org.ID, ActionType: models.ApprovalActionPayout,
		ResourceID: "p1", ResourceType: "payout", RequestedByID: director.ID,
		RequiredCount: 1, Payload: []byte(`{}`),
	})
	svc.CreateRequest(context.Background(), services.CreateApprovalInput{
		OrganizationID: org.ID, ActionType: models.ApprovalActionConversion,
		ResourceID: "c1", ResourceType: "conversion", RequestedByID: director.ID,
		RequiredCount: 1, Payload: []byte(`{}`),
	})

	pending, err := svc.GetPendingRequests(org.ID)
	require.NoError(t, err)
	assert.Len(t, pending, 2)
}
