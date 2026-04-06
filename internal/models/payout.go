package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayoutType string

const (
	PayoutTypeBank   PayoutType = "bank"
	PayoutTypeCrypto PayoutType = "crypto"
)

type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
	PayoutStatusCancelled  PayoutStatus = "cancelled"
)

type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "pending"
	BatchStatusProcessing BatchStatus = "processing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusFailed     BatchStatus = "failed"
	BatchStatusCancelled  BatchStatus = "cancelled"
)

// Payout represents a payout transaction
type Payout struct {
	ID             string          `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Type           PayoutType      `gorm:"type:text;not null" json:"type"`
	Amount         decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency       string          `gorm:"type:char(3);not null" json:"currency"`
	BeneficiaryID  string          `gorm:"type:uuid;not null" json:"beneficiaryId"`
	Reference      *string         `json:"reference"`
	Description    *string         `json:"description"`
	Status         PayoutStatus    `gorm:"type:text;default:'pending'" json:"status"`
	ProviderRef    *string         `gorm:"type:text" json:"-"`
	InitiatedByID  *string         `gorm:"type:uuid" json:"initiatedById"`
	Fee            decimal.Decimal `gorm:"type:decimal(20,8);default:0" json:"fee"`
	ExecutedAt     *time.Time      `json:"executedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Beneficiary    Beneficiary     `gorm:"foreignKey:BeneficiaryID;references:ID" json:"-"`
}

func (Payout) TableName() string {
	return "payouts"
}

// PayoutBatch represents a batch of payouts
type PayoutBatch struct {
	ID             string          `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Name           string          `gorm:"not null" json:"name"`
	TotalAmount    decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"totalAmount"`
	TotalCount     int             `gorm:"not null" json:"totalCount"`
	ProcessedCount int             `gorm:"default:0" json:"processedCount"`
	Status         BatchStatus     `gorm:"type:text;default:'pending'" json:"status"`
	FileURL        *string         `json:"fileUrl"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (PayoutBatch) TableName() string {
	return "payout_batches"
}
