package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type ConversionType string

const (
	ConversionTypeMarket ConversionType = "market"
	ConversionTypeLimit  ConversionType = "limit"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// ConversionOrder represents a currency conversion order
type ConversionOrder struct {
	ID             string           `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string           `gorm:"type:uuid;not null" json:"organizationId"`
	FromCurrency   string           `gorm:"type:char(3);not null" json:"fromCurrency"`
	ToCurrency     string           `gorm:"type:char(3);not null" json:"toCurrency"`
	FromAmount     decimal.Decimal  `gorm:"type:decimal(20,8);not null" json:"fromAmount"`
	ToAmount       decimal.Decimal  `gorm:"type:decimal(20,8);not null" json:"toAmount"`
	Rate           decimal.Decimal  `gorm:"type:decimal(20,8);not null" json:"rate"`
	Spread         decimal.Decimal  `gorm:"type:decimal(20,8);default:0" json:"spread"`
	Fee            decimal.Decimal  `gorm:"type:decimal(20,8);default:0" json:"fee"`
	ProviderRef    *string          `gorm:"type:text" json:"-"`
	InitiatedByID  *string          `gorm:"type:uuid" json:"initiatedById"`
	Status         OrderStatus      `gorm:"type:text;default:'pending'" json:"status"`
	Type           ConversionType   `gorm:"type:text;default:'market'" json:"type"`
	LimitPrice     *decimal.Decimal `gorm:"type:decimal(20,8)" json:"limitPrice"`
	ExecutedAt     *time.Time       `json:"executedAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Organization   Organization     `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ConversionOrder) TableName() string {
	return "conversion_orders"
}
