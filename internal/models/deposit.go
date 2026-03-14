package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type DepositType string

const (
	DepositTypeFiat   DepositType = "fiat"
	DepositTypeCrypto DepositType = "crypto"
)

type DepositStatus string

const (
	DepositStatusPending    DepositStatus = "pending"
	DepositStatusProcessing DepositStatus = "processing"
	DepositStatusCompleted  DepositStatus = "completed"
	DepositStatusFailed     DepositStatus = "failed"
	DepositStatusCancelled  DepositStatus = "cancelled"
)

// Deposit represents a wallet deposit
type Deposit struct {
	ID                string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID    string          `gorm:"type:uuid;not null" json:"organizationId"`
	WalletID          string          `gorm:"type:uuid;not null" json:"walletId"`
	Type              DepositType     `gorm:"type:text;not null" json:"type"`
	Amount            decimal.Decimal `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency          string          `gorm:"type:char(3);not null" json:"currency"`
	Reference         *string         `json:"reference"`
	Status            DepositStatus   `gorm:"type:text;default:'pending'" json:"status"`
	ProviderReference *string         `json:"providerReference"`
	ExecutedAt        *time.Time      `json:"executedAt"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	Organization      Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (Deposit) TableName() string {
	return "deposits"
}
