package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type WalletType string

const (
	WalletTypeFiat   WalletType = "fiat"
	WalletTypeCrypto WalletType = "crypto"
)

type TransactionType string

const (
	TransactionTypeDeposit     TransactionType = "deposit"
	TransactionTypeWithdrawal  TransactionType = "withdrawal"
	TransactionTypeTransferIn  TransactionType = "transfer_in"
	TransactionTypeTransferOut TransactionType = "transfer_out"
	TransactionTypeConversion  TransactionType = "conversion"
	TransactionTypeFee         TransactionType = "fee"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusProcessing TransactionStatus = "processing"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
)

// Wallet represents a currency wallet
type Wallet struct {
	ID               string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID   string              `gorm:"type:uuid;not null" json:"organizationId"`
	Type             WalletType          `gorm:"type:text;not null" json:"type"`
	Currency         string              `gorm:"type:char(3);not null" json:"currency"`
	Balance          decimal.Decimal     `gorm:"type:decimal(20,8);default:0" json:"balance"`
	AvailableBalance decimal.Decimal     `gorm:"type:decimal(20,8);default:0" json:"availableBalance"`
	PendingBalance   decimal.Decimal     `gorm:"type:decimal(20,8);default:0" json:"pendingBalance"`
	IsActive         bool                `gorm:"default:true" json:"isActive"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
	Organization     Organization        `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Transactions     []WalletTransaction `gorm:"foreignKey:WalletID" json:"-"`
}

func (Wallet) TableName() string {
	return "wallets"
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID          string            `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WalletID    string            `gorm:"type:uuid;not null" json:"walletId"`
	Type        TransactionType   `gorm:"type:text;not null" json:"type"`
	Amount      decimal.Decimal   `gorm:"type:decimal(20,8);not null" json:"amount"`
	Currency    string            `gorm:"type:char(3);not null" json:"currency"`
	Description string            `gorm:"not null" json:"description"`
	Reference   *string           `json:"reference"`
	Status      TransactionStatus `gorm:"type:text;default:'pending'" json:"status"`
	Metadata    datatypes.JSON    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Wallet      Wallet            `gorm:"foreignKey:WalletID;references:ID" json:"-"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
