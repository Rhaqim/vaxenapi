package services

import (
	"errors"

	"vaxen/api/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

// --- Platform Settings ---

func (s *AdminService) GetSettings(category string) ([]models.PlatformSetting, error) {
	var settings []models.PlatformSetting
	q := s.db
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if err := q.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *AdminService) UpsertSetting(key, value, category, updatedBy string) error {
	var setting models.PlatformSetting
	err := s.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return s.db.Create(&models.PlatformSetting{
			Key:       key,
			Value:     value,
			Category:  category,
			UpdatedBy: updatedBy,
		}).Error
	}
	return s.db.Model(&setting).Updates(map[string]any{
		"value":      value,
		"category":   category,
		"updated_by": updatedBy,
	}).Error
}

// --- Feature Flags ---

func (s *AdminService) GetFeatureFlags() ([]models.FeatureFlag, error) {
	var flags []models.FeatureFlag
	if err := s.db.Find(&flags).Error; err != nil {
		return nil, err
	}
	return flags, nil
}

func (s *AdminService) SetFeatureFlag(name string, enabled bool, description string, updatedBy string) error {
	var flag models.FeatureFlag
	err := s.db.Where("name = ?", name).First(&flag).Error
	if err != nil {
		return s.db.Create(&models.FeatureFlag{
			Name:        name,
			Enabled:     enabled,
			Description: description,
			UpdatedBy:   updatedBy,
		}).Error
	}
	return s.db.Model(&flag).Updates(map[string]any{
		"enabled":     enabled,
		"description": description,
		"updated_by":  updatedBy,
	}).Error
}

func (s *AdminService) IsFeatureEnabled(name string) bool {
	var flag models.FeatureFlag
	if err := s.db.Where("name = ? AND enabled = ?", name, true).First(&flag).Error; err != nil {
		return false
	}
	return true
}

// --- Exchange Rates ---

func (s *AdminService) GetExchangeRates() ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	if err := s.db.Where("is_active = ?", true).Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (s *AdminService) UpsertExchangeRate(from, to string, rate, spread decimal.Decimal, updatedBy string) error {
	var existing models.ExchangeRate
	err := s.db.Where("from_currency = ? AND to_currency = ?", from, to).First(&existing).Error
	if err != nil {
		return s.db.Create(&models.ExchangeRate{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         rate,
			Spread:       spread,
			IsActive:     true,
			UpdatedBy:    updatedBy,
		}).Error
	}
	return s.db.Model(&existing).Updates(map[string]any{
		"rate":       rate,
		"spread":     spread,
		"updated_by": updatedBy,
	}).Error
}

// --- User & Organization Management ---

func (s *AdminService) GetAllUsers(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	s.db.Model(&models.User{}).Count(&total)
	if err := s.db.Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *AdminService) GetAllOrganizations(page, limit int) ([]models.Organization, int64, error) {
	var orgs []models.Organization
	var total int64
	s.db.Model(&models.Organization{}).Count(&total)
	if err := s.db.Offset((page - 1) * limit).Limit(limit).Find(&orgs).Error; err != nil {
		return nil, 0, err
	}
	return orgs, total, nil
}

func (s *AdminService) ApproveOrganization(orgID string) error {
	var org models.Organization
	if err := s.db.First(&org, "id = ?", orgID).Error; err != nil {
		return errors.New("organization not found")
	}
	return s.db.Model(&org).Update("kyb_status", models.KybStatusApproved).Error
}

func (s *AdminService) RejectOrganization(orgID string, reason string) error {
	var org models.Organization
	if err := s.db.First(&org, "id = ?", orgID).Error; err != nil {
		return errors.New("organization not found")
	}
	return s.db.Model(&org).Update("kyb_status", models.KybStatusRejected).Error
}
