package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// PlatformSetting stores global platform configuration managed by admins.
type PlatformSetting struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	Category  string    `gorm:"type:text;not null;default:'general'" json:"category"` // general, exchange, features, providers
	UpdatedBy string    `gorm:"type:uuid" json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (PlatformSetting) TableName() string { return "platform_settings" }

// FeatureFlag controls feature availability on the platform.
type FeatureFlag struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Enabled     bool      `gorm:"default:false" json:"enabled"`
	Description string    `json:"description"`
	UpdatedBy   string    `gorm:"type:uuid" json:"updatedBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (FeatureFlag) TableName() string { return "feature_flags" }

// ExchangeRate stores admin-managed exchange rates (for the internal exchange provider).
type ExchangeRate struct {
	ID           string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FromCurrency string          `gorm:"type:char(5);not null" json:"fromCurrency"`
	ToCurrency   string          `gorm:"type:char(5);not null" json:"toCurrency"`
	Rate         decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"rate"`
	Spread       decimal.Decimal `gorm:"type:decimal(10,6);default:0" json:"spread"`
	IsActive     bool            `gorm:"default:true" json:"isActive"`
	UpdatedBy    string          `gorm:"type:uuid" json:"updatedBy"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

func (ExchangeRate) TableName() string { return "exchange_rates" }
