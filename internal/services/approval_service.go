package services

import (
	"context"
	"errors"
	"time"

	"vaxen/api/internal/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ApprovalService struct {
	db *gorm.DB
}

func NewApprovalService(db *gorm.DB) *ApprovalService {
	return &ApprovalService{db: db}
}

type CreateApprovalInput struct {
	OrganizationID string
	ActionType     models.ApprovalActionType
	ResourceID     string
	ResourceType   string
	RequestedByID  string
	RequiredCount  int
	Payload        []byte
}

// GetRequiredApprovals returns the number of approvals required for an action type.
func (s *ApprovalService) GetRequiredApprovals(orgID string, actionType models.ApprovalActionType) (int, error) {
	var policy models.ApprovalPolicy
	if err := s.db.Where("organization_id = ? AND action_type = ?", orgID, actionType).First(&policy).Error; err != nil {
		return 1, nil // Default to 1 if no policy exists
	}
	return policy.RequiredApprovals, nil
}

// CreateRequest creates a new approval request.
func (s *ApprovalService) CreateRequest(_ context.Context, input CreateApprovalInput) (*models.ApprovalRequest, error) {
	expiry := time.Now().Add(72 * time.Hour)
	req := models.ApprovalRequest{
		OrganizationID: input.OrganizationID,
		ActionType:     input.ActionType,
		ResourceID:     input.ResourceID,
		ResourceType:   input.ResourceType,
		RequestedByID:  input.RequestedByID,
		Status:         models.ApprovalStatusPending,
		RequiredCount:  input.RequiredCount,
		CurrentCount:   0,
		Payload:        datatypes.JSON(input.Payload),
		ExpiresAt:      &expiry,
	}
	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// Vote submits a director's vote on an approval request.
func (s *ApprovalService) Vote(_ context.Context, requestID string, voterID string, decision models.ApprovalStatus, comment string) error {
	var req models.ApprovalRequest
	if err := s.db.First(&req, "id = ?", requestID).Error; err != nil {
		return errors.New("approval request not found")
	}

	if req.Status != models.ApprovalStatusPending {
		return errors.New("approval request is no longer pending")
	}

	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		s.db.Model(&req).Update("status", models.ApprovalStatusExpired)
		return errors.New("approval request has expired")
	}

	// Check voter hasn't already voted
	var existingVote models.ApprovalVote
	if err := s.db.Where("approval_request_id = ? AND voter_id = ?", requestID, voterID).First(&existingVote).Error; err == nil {
		return errors.New("you have already voted on this request")
	}

	// Verify voter is a director in the same org
	var voter models.User
	if err := s.db.First(&voter, "id = ? AND organization_id = ?", voterID, req.OrganizationID).Error; err != nil {
		return errors.New("voter not found in organization")
	}
	if !voter.Role.CanApprove() {
		return errors.New("insufficient permissions to approve")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		vote := models.ApprovalVote{
			ApprovalRequestID: requestID,
			VoterID:           voterID,
			Decision:          decision,
			Comment:           comment,
		}
		if err := tx.Create(&vote).Error; err != nil {
			return err
		}

		if decision == models.ApprovalStatusRejected {
			now := time.Now()
			return tx.Model(&req).Updates(map[string]any{
				"status":       models.ApprovalStatusRejected,
				"completed_at": now,
			}).Error
		}

		// Count approvals
		newCount := req.CurrentCount + 1
		updates := map[string]any{"current_count": newCount}

		if newCount >= req.RequiredCount {
			now := time.Now()
			updates["status"] = models.ApprovalStatusApproved
			updates["completed_at"] = now
		}

		return tx.Model(&req).Updates(updates).Error
	})
}

// GetPendingRequests returns all pending approval requests for an organization.
func (s *ApprovalService) GetPendingRequests(orgID string) ([]models.ApprovalRequest, error) {
	var requests []models.ApprovalRequest
	if err := s.db.Preload("Votes").Where("organization_id = ? AND status = ?", orgID, models.ApprovalStatusPending).
		Order("created_at desc").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetRequest returns a specific approval request.
func (s *ApprovalService) GetRequest(requestID string, orgID string) (*models.ApprovalRequest, error) {
	var req models.ApprovalRequest
	if err := s.db.Preload("Votes").Where("id = ? AND organization_id = ?", requestID, orgID).First(&req).Error; err != nil {
		return nil, errors.New("approval request not found")
	}
	return &req, nil
}

// UpdatePolicy updates the approval policy for an organization.
func (s *ApprovalService) UpdatePolicy(orgID string, actionType models.ApprovalActionType, requiredApprovals int) error {
	var policy models.ApprovalPolicy
	err := s.db.Where("organization_id = ? AND action_type = ?", orgID, actionType).First(&policy).Error
	if err != nil {
		// Create new policy
		policy = models.ApprovalPolicy{
			OrganizationID:    orgID,
			ActionType:        actionType,
			RequiredApprovals: requiredApprovals,
		}
		return s.db.Create(&policy).Error
	}
	return s.db.Model(&policy).Update("required_approvals", requiredApprovals).Error
}
