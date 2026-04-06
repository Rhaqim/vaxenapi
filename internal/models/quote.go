package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type QuoteType string

const (
	QuoteTypeSpot        QuoteType = "spot"
	QuoteTypeAutoConvert QuoteType = "auto_convert"
	QuoteTypeLimit       QuoteType = "limit"
)

// Quote represents a currency conversion quote
type Quote struct {
	ID           string          `gorm:"type:uuid;primary_key" json:"id"`
	FromCurrency string          `gorm:"type:char(3);not null" json:"fromCurrency"`
	ToCurrency   string          `gorm:"type:char(3);not null" json:"toCurrency"`
	FromAmount   decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"fromAmount"`
	ToAmount     decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"toAmount"`
	Rate         decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"rate"`
	Spread       decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"spread"`
	Fee          decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"fee"`
	FeeCurrency  string          `gorm:"type:char(3);not null" json:"feeCurrency"`
	ExpiresAt    time.Time       `json:"expiresAt"`
	Type         QuoteType       `gorm:"type:text;not null" json:"type"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

func (Quote) TableName() string {
	return "quotes"
}
