package models

import "time"

// Web3Wallet represents a blockchain wallet managed by the platform.
// Private keys are stored in the configured key management provider (e.g. AWS KMS).
type Web3Wallet struct {
	ID             string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string       `gorm:"type:uuid;not null" json:"organizationId"`
	Address        string       `gorm:"type:text;not null" json:"address"`
	Network        string       `gorm:"type:text;not null" json:"network"` // "ethereum", "polygon", "bitcoin"
	KeyID          string       `gorm:"type:text;not null" json:"-"`       // Provider key reference (e.g. KMS ARN)
	Label          string       `gorm:"type:text" json:"label"`
	IsActive       bool         `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (Web3Wallet) TableName() string { return "web3_wallets" }
