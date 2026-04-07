package services

import (
	"errors"

	"vaxen/api/internal/models"

	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

type CreateAccountInput struct {
	OrganizationID string `json:"-"`
	Name           string `json:"name" binding:"required"`
	Currency       string `json:"currency" binding:"required"`
	Type           string `json:"type" binding:"required"` // iban, pix, ach, swift
	AccountNumber  string `json:"accountNumber" binding:"required"`
	RoutingNumber  string `json:"routingNumber"`
	BankCode       string `json:"bankCode"`
	BankName       string `json:"bankName" binding:"required"`
	BankCountry    string `json:"bankCountry" binding:"required"`
}

func (s *AccountService) List(orgID string) ([]models.AccountNumber, error) {
	var accounts []models.AccountNumber
	if err := s.db.Where("organization_id = ? AND is_active = ?", orgID, true).
		Order("created_at desc").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *AccountService) Get(orgID, id string) (*models.AccountNumber, error) {
	var a models.AccountNumber
	if err := s.db.Where("id = ? AND organization_id = ?", id, orgID).First(&a).Error; err != nil {
		return nil, errors.New("account not found")
	}
	return &a, nil
}

func (s *AccountService) Create(input CreateAccountInput) (*models.AccountNumber, error) {
	aType := models.AccountType(input.Type)
	if aType != models.AccountTypeIBAN && aType != models.AccountTypePIX &&
		aType != models.AccountTypeACH && aType != models.AccountTypeSWIFT {
		return nil, errors.New("type must be 'iban', 'pix', 'ach', or 'swift'")
	}

	a := models.AccountNumber{
		OrganizationID: input.OrganizationID,
		Name:           input.Name,
		Currency:       input.Currency,
		Type:           aType,
		AccountNumber:  input.AccountNumber,
		BankName:       input.BankName,
		BankCountry:    input.BankCountry,
		IsActive:       true,
	}

	if input.RoutingNumber != "" {
		a.RoutingNumber = &input.RoutingNumber
	}
	if input.BankCode != "" {
		a.BankCode = &input.BankCode
	}

	if err := s.db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
