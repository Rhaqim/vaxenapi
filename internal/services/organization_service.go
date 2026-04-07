package services

import (
	"errors"

	"vaxen/api/internal/models"

	"gorm.io/gorm"
)

type OrganizationService struct {
	db *gorm.DB
}

func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

type UpdateOrganizationInput struct {
	Name      *string `json:"name"`
	LegalName *string `json:"legalName"`
	TaxID     *string `json:"taxId"`
	Country   *string `json:"country"`
}

// Get returns the organization for the authenticated user.
func (s *OrganizationService) Get(orgID string) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.First(&org, "id = ?", orgID).Error; err != nil {
		return nil, errors.New("organization not found")
	}
	return &org, nil
}

// Update updates organization fields.
func (s *OrganizationService) Update(orgID string, input UpdateOrganizationInput) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.First(&org, "id = ?", orgID).Error; err != nil {
		return nil, errors.New("organization not found")
	}

	updates := make(map[string]any)
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.LegalName != nil {
		updates["legal_name"] = *input.LegalName
	}
	if input.TaxID != nil {
		updates["tax_id"] = *input.TaxID
	}
	if input.Country != nil {
		updates["country"] = *input.Country
	}

	if len(updates) > 0 {
		if err := s.db.Model(&org).Updates(updates).Error; err != nil {
			return nil, err
		}
		s.db.First(&org, "id = ?", orgID)
	}

	return &org, nil
}

// GetUsers returns all users in the organization.
func (s *OrganizationService) GetUsers(orgID string) ([]models.User, error) {
	var users []models.User
	if err := s.db.Where("organization_id = ? AND is_active = ?", orgID, true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
