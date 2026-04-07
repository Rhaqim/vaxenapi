package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

// SeedExchangeRatesInput represents a request to bulk-seed exchange rates.
type SeedExchangeRatesInput struct {
	Base          string   `json:"base" binding:"required"`          // e.g. "USD"
	Targets       []string `json:"targets" binding:"required,min=1"` // e.g. ["EUR", "GBP"]
	DefaultSpread string   `json:"defaultSpread"`                    // e.g. "0.005", defaults to 0.5%
}

// SeedExchangeRateResult is returned for each pair processed.
type SeedExchangeRateResult struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Rate   string `json:"rate"`
	Status string `json:"status"` // "created", "updated", "failed"
	Error  string `json:"error,omitempty"`
}

// SeedExchangeRates fetches live rates from a public API and upserts them.
func (s *AdminService) SeedExchangeRates(input SeedExchangeRatesInput, updatedBy string) ([]SeedExchangeRateResult, error) {
	base := strings.ToUpper(input.Base)
	targets := make([]string, len(input.Targets))
	for i, t := range input.Targets {
		targets[i] = strings.ToUpper(t)
	}

	spread := decimal.NewFromFloat(0.005)
	if input.DefaultSpread != "" {
		if s, err := decimal.NewFromString(input.DefaultSpread); err == nil {
			spread = s
		}
	}

	rates, err := fetchPublicRates(base, targets)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}

	results := make([]SeedExchangeRateResult, 0, len(targets))
	for _, target := range targets {
		rate, ok := rates[target]
		if !ok {
			results = append(results, SeedExchangeRateResult{
				From: base, To: target, Status: "failed",
				Error: "rate not found in API response",
			})
			continue
		}

		err := s.UpsertExchangeRate(base, target, rate, spread, updatedBy)
		status := "updated"
		errMsg := ""
		if err != nil {
			status = "failed"
			errMsg = err.Error()
		} else {
			// Check if it was a create or update
			var count int64
			s.db.Model(&models.ExchangeRate{}).Where("from_currency = ? AND to_currency = ?", base, target).Count(&count)
			if count == 1 {
				status = "created"
			}
		}
		results = append(results, SeedExchangeRateResult{
			From: base, To: target, Rate: rate.StringFixed(6), Status: status, Error: errMsg,
		})
	}

	return results, nil
}

// fetchPublicRates fetches exchange rates from the Frankfurter API (ECB data, no key required).
func fetchPublicRates(base string, targets []string) (map[string]decimal.Decimal, error) {
	url := fmt.Sprintf("https://api.frankfurter.app/latest?from=%s&to=%s",
		base, strings.Join(targets, ","))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Base  string             `json:"base"`
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	rates := make(map[string]decimal.Decimal, len(result.Rates))
	for currency, rate := range result.Rates {
		rates[currency] = decimal.NewFromFloat(rate)
	}

	return rates, nil
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
