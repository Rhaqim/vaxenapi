package services

import (
	"errors"

	"vaxen/api/internal/models"
	"vaxen/api/internal/providers"

	"gorm.io/gorm"
)

type BeneficiaryService struct {
	db  *gorm.DB
	reg *providers.Registry
}

func NewBeneficiaryService(db *gorm.DB, reg *providers.Registry) *BeneficiaryService {
	return &BeneficiaryService{db: db, reg: reg}
}

type CreateBeneficiaryInput struct {
	OrganizationID string `json:"-"`
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"` // "bank" or "crypto"
	Currency       string `json:"currency" binding:"required"`
	AccountNumber  string `json:"accountNumber"`
	RoutingNumber  string `json:"routingNumber"`
	BankCode       string `json:"bankCode"`
	BankName       string `json:"bankName"`
	BankCountry    string `json:"bankCountry"`
	Address        string `json:"address"` // Crypto wallet address or bank address
	Network        string `json:"network"` // For crypto beneficiaries
}

type UpdateBeneficiaryInput struct {
	Name          *string `json:"name"`
	AccountNumber *string `json:"accountNumber"`
	RoutingNumber *string `json:"routingNumber"`
	BankCode      *string `json:"bankCode"`
	BankName      *string `json:"bankName"`
	BankCountry   *string `json:"bankCountry"`
	Address       *string `json:"address"`
	Network       *string `json:"network"`
	Currency      *string `json:"currency"`
}

func (s *BeneficiaryService) List(orgID string) ([]models.Beneficiary, error) {
	var beneficiaries []models.Beneficiary
	if err := s.db.Where("organization_id = ? AND is_active = ?", orgID, true).
		Order("created_at desc").Find(&beneficiaries).Error; err != nil {
		return nil, err
	}
	return beneficiaries, nil
}

func (s *BeneficiaryService) Get(orgID, id string) (*models.Beneficiary, error) {
	var b models.Beneficiary
	if err := s.db.Where("id = ? AND organization_id = ?", id, orgID).First(&b).Error; err != nil {
		return nil, errors.New("beneficiary not found")
	}
	return &b, nil
}

func (s *BeneficiaryService) Create(input CreateBeneficiaryInput) (*models.Beneficiary, error) {
	bType := models.BeneficiaryType(input.Type)
	if bType != models.BeneficiaryTypeBank && bType != models.BeneficiaryTypeCrypto {
		return nil, errors.New("type must be 'bank' or 'crypto'")
	}

	b := models.Beneficiary{
		OrganizationID: input.OrganizationID,
		Name:           input.Name,
		Type:           bType,
		Currency:       input.Currency,
		IsActive:       true,
	}

	if input.AccountNumber != "" {
		b.AccountNumber = &input.AccountNumber
	}
	if input.RoutingNumber != "" {
		b.RoutingNumber = &input.RoutingNumber
	}
	if input.BankCode != "" {
		b.BankCode = &input.BankCode
	}
	if input.BankName != "" {
		b.BankName = &input.BankName
	}
	if input.BankCountry != "" {
		b.BankCountry = &input.BankCountry
	}
	if input.Address != "" {
		b.Address = &input.Address
	}
	if input.Network != "" {
		b.Network = &input.Network
	}

	if err := s.db.Create(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *BeneficiaryService) Update(orgID, id string, input UpdateBeneficiaryInput) (*models.Beneficiary, error) {
	var b models.Beneficiary
	if err := s.db.Where("id = ? AND organization_id = ?", id, orgID).First(&b).Error; err != nil {
		return nil, errors.New("beneficiary not found")
	}

	updates := make(map[string]any)
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.AccountNumber != nil {
		updates["account_number"] = *input.AccountNumber
	}
	if input.RoutingNumber != nil {
		updates["routing_number"] = *input.RoutingNumber
	}
	if input.BankCode != nil {
		updates["bank_code"] = *input.BankCode
	}
	if input.BankName != nil {
		updates["bank_name"] = *input.BankName
	}
	if input.BankCountry != nil {
		updates["bank_country"] = *input.BankCountry
	}
	if input.Address != nil {
		updates["address"] = *input.Address
	}
	if input.Network != nil {
		updates["network"] = *input.Network
	}
	if input.Currency != nil {
		updates["currency"] = *input.Currency
	}

	if len(updates) == 0 {
		return &b, nil
	}

	if err := s.db.Model(&b).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.db.First(&b, "id = ?", id)
	return &b, nil
}

func (s *BeneficiaryService) Delete(orgID, id string) error {
	var b models.Beneficiary
	if err := s.db.Where("id = ? AND organization_id = ?", id, orgID).First(&b).Error; err != nil {
		return errors.New("beneficiary not found")
	}

	// Soft-delete: mark inactive rather than removing
	return s.db.Model(&b).Update("is_active", false).Error
}
