package services

import (
	"context"
	"errors"
	"time"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/kyc"

	"gorm.io/gorm"
)

type KYCService struct {
	db  *gorm.DB
	reg *providers.Registry
}

func NewKYCService(db *gorm.DB, reg *providers.Registry) *KYCService {
	return &KYCService{db: db, reg: reg}
}

type KYBSubmitInput struct {
	OrganizationID     string             `json:"-"` // set from context
	LegalName          string             `json:"legalName" binding:"required"`
	RegistrationNumber string             `json:"registrationNumber" binding:"required"`
	TaxID              string             `json:"taxId"`
	Country            string             `json:"country" binding:"required"`
	Address            kyc.Address        `json:"address" binding:"required"`
	Directors          []kyc.DirectorInfo `json:"directors" binding:"required"`
	Documents          []kyc.Document     `json:"documents,omitempty"`
}

// SubmitKYB sends KYB data to the provider and updates the org record.
func (s *KYCService) SubmitKYB(ctx context.Context, input KYBSubmitInput) (*kyc.KYBResult, error) {
	provider, err := s.reg.RequireKYC()
	if err != nil {
		return nil, err
	}

	var org models.Organization
	if err := s.db.First(&org, "id = ?", input.OrganizationID).Error; err != nil {
		return nil, errors.New("organization not found")
	}

	if org.KYBStatus == models.KybStatusApproved {
		return nil, errors.New("KYB already approved")
	}

	result, err := provider.SubmitKYB(ctx, kyc.KYBSubmission{
		OrganizationID:     input.OrganizationID,
		LegalName:          input.LegalName,
		RegistrationNumber: input.RegistrationNumber,
		TaxID:              input.TaxID,
		Country:            input.Country,
		Address:            input.Address,
		Directors:          input.Directors,
		Documents:          input.Documents,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	s.db.Model(&org).Updates(map[string]any{
		"kyb_status":       models.KybStatusPending,
		"kyb_reference_id": result.ReferenceID,
		"kyb_submitted_at": now,
	})

	return result, nil
}

// GetKYBStatus checks the provider for the latest KYB status and updates our record.
func (s *KYCService) GetKYBStatus(ctx context.Context, organizationID string) (*kyc.KYBResult, error) {
	provider, err := s.reg.RequireKYC()
	if err != nil {
		return nil, err
	}

	var org models.Organization
	if err := s.db.First(&org, "id = ?", organizationID).Error; err != nil {
		return nil, errors.New("organization not found")
	}

	if org.KYBReferenceID == "" {
		return &kyc.KYBResult{
			Status: kyc.VerificationStatus(org.KYBStatus),
		}, nil
	}

	result, err := provider.GetKYBStatus(ctx, org.KYBReferenceID)
	if err != nil {
		return nil, err
	}

	newStatus := mapVerificationStatus(result.Status)
	updates := map[string]any{"kyb_status": newStatus}
	if newStatus == models.KybStatusApproved {
		now := time.Now()
		updates["kyb_approved_at"] = now
	}
	s.db.Model(&org).Updates(updates)

	return result, nil
}

// SubmitKYC sends director identity verification data to the provider.
func (s *KYCService) SubmitKYC(ctx context.Context, input kyc.KYCSubmission) (*kyc.KYCResult, error) {
	provider, err := s.reg.RequireKYC()
	if err != nil {
		return nil, err
	}
	return provider.SubmitKYC(ctx, input)
}

// GetKYCStatus checks director KYC status from the provider.
func (s *KYCService) GetKYCStatus(ctx context.Context, referenceID string) (*kyc.KYCResult, error) {
	provider, err := s.reg.RequireKYC()
	if err != nil {
		return nil, err
	}
	return provider.GetKYCStatus(ctx, referenceID)
}

// GenerateVerificationURL creates a hosted verification link for a director.
func (s *KYCService) GenerateVerificationURL(ctx context.Context, applicantID string, redirectURL string) (string, error) {
	provider, err := s.reg.RequireKYC()
	if err != nil {
		return "", err
	}
	return provider.GenerateVerificationURL(ctx, applicantID, redirectURL)
}

func mapVerificationStatus(status kyc.VerificationStatus) models.KybStatus {
	switch status {
	case kyc.StatusApproved:
		return models.KybStatusApproved
	case kyc.StatusRejected:
		return models.KybStatusRejected
	case kyc.StatusRequiresInfo:
		return models.KybStatusRequiresInfo
	default:
		return models.KybStatusPending
	}
}
