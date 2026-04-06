package services

import (
	"context"
	"encoding/json"
	"errors"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/payment"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentService struct {
	db       *gorm.DB
	reg      *providers.Registry
	approval *ApprovalService
}

func NewPaymentService(db *gorm.DB, reg *providers.Registry, approval *ApprovalService) *PaymentService {
	return &PaymentService{db: db, reg: reg, approval: approval}
}

type SendPaymentInput struct {
	OrganizationID string          `json:"-"` // set from context
	InitiatedByID  string          `json:"-"` // set from context
	BeneficiaryID  string          `json:"beneficiaryId" binding:"required"`
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	Currency       string          `json:"currency" binding:"required"`
	Reference      string          `json:"reference"`
	Description    string          `json:"description"`
}

// InitiatePayment creates a payout record and submits it for approval if required.
func (s *PaymentService) InitiatePayment(ctx context.Context, input SendPaymentInput) (*models.Payout, error) {
	var beneficiary models.Beneficiary
	if err := s.db.Where("id = ? AND organization_id = ?", input.BeneficiaryID, input.OrganizationID).First(&beneficiary).Error; err != nil {
		return nil, errors.New("beneficiary not found")
	}

	payout := models.Payout{
		OrganizationID: input.OrganizationID,
		Type:           models.PayoutTypeBank,
		Amount:         input.Amount,
		Currency:       input.Currency,
		BeneficiaryID:  input.BeneficiaryID,
		Reference:      &input.Reference,
		Description:    &input.Description,
		InitiatedByID:  &input.InitiatedByID,
		Status:         models.PayoutStatusPending,
	}
	if err := s.db.Create(&payout).Error; err != nil {
		return nil, err
	}

	required, err := s.approval.GetRequiredApprovals(input.OrganizationID, models.ApprovalActionPayout)
	if err != nil {
		return nil, err
	}

	if required > 1 {
		payload, _ := json.Marshal(input)
		if _, err := s.approval.CreateRequest(ctx, CreateApprovalInput{
			OrganizationID: input.OrganizationID,
			ActionType:     models.ApprovalActionPayout,
			ResourceID:     payout.ID,
			ResourceType:   "payout",
			RequestedByID:  input.InitiatedByID,
			RequiredCount:  required,
			Payload:        payload,
		}); err != nil {
			return nil, err
		}
		return &payout, nil
	}

	if err := s.ExecutePayment(ctx, &payout, &beneficiary); err != nil {
		return nil, err
	}

	return &payout, nil
}

// ExecutePayment sends the payment via the configured provider.
func (s *PaymentService) ExecutePayment(ctx context.Context, payout *models.Payout, beneficiary *models.Beneficiary) error {
	provider, err := s.reg.RequirePayment()
	if err != nil {
		return err
	}

	ref := ""
	if payout.Reference != nil {
		ref = *payout.Reference
	}
	desc := ""
	if payout.Description != nil {
		desc = *payout.Description
	}

	result, err := provider.SendPayment(ctx, payment.SendRequest{
		OrganizationID:  payout.OrganizationID,
		Amount:          payout.Amount,
		Currency:        payout.Currency,
		Reference:       ref,
		Description:     desc,
		BeneficiaryID:   beneficiary.ID,
		BeneficiaryName: beneficiary.Name,
	})
	if err != nil {
		s.db.Model(payout).Update("status", models.PayoutStatusFailed)
		return err
	}

	s.db.Model(payout).Updates(map[string]any{
		"status":       models.PayoutStatusProcessing,
		"provider_ref": result.ProviderRef,
		"fee":          result.Fee,
	})

	return nil
}

// GetPaymentStatus checks the provider for the latest status.
func (s *PaymentService) GetPaymentStatus(ctx context.Context, payoutID string, orgID string) (*models.Payout, error) {
	var payout models.Payout
	if err := s.db.Where("id = ? AND organization_id = ?", payoutID, orgID).First(&payout).Error; err != nil {
		return nil, errors.New("payout not found")
	}

	if payout.ProviderRef != nil {
		if provider, err := s.reg.RequirePayment(); err == nil {
			if status, err := provider.GetPaymentStatus(ctx, *payout.ProviderRef); err == nil {
				s.db.Model(&payout).Update("status", status.Status)
				payout.Status = models.PayoutStatus(status.Status)
			}
		}
	}

	return &payout, nil
}

// GetPayouts returns all payouts for an organization.
func (s *PaymentService) GetPayouts(orgID string) ([]models.Payout, error) {
	var payouts []models.Payout
	if err := s.db.Where("organization_id = ?", orgID).Order("created_at desc").Find(&payouts).Error; err != nil {
		return nil, err
	}
	return payouts, nil
}

// GetFees estimates fees for a payment.
func (s *PaymentService) GetFees(ctx context.Context, amount decimal.Decimal, currency string, country string, method string) (*payment.FeeResult, error) {
	provider, err := s.reg.RequirePayment()
	if err != nil {
		return nil, err
	}
	return provider.GetFees(ctx, payment.FeeRequest{
		Amount:   amount,
		Currency: currency,
		Country:  country,
		Method:   method,
	})
}
