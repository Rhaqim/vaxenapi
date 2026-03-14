package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// LimitOrder represents a limit order
type LimitOrder struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	FromCurrency   string          `gorm:"type:char(3);not null" json:"fromCurrency"`
	ToCurrency     string          `gorm:"type:char(3);not null" json:"toCurrency"`
	Amount         decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`
	LimitPrice     decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"limitPrice"`
	Type           OrderType       `gorm:"type:text;default:'limit'" json:"type"`
	Status         OrderStatus     `gorm:"type:text;default:'pending'" json:"status"`
	ExecutedAt     *time.Time      `json:"executedAt"`
	CancelledAt    *time.Time      `json:"cancelledAt"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (LimitOrder) TableName() string {
	return "limit_orders"
}
