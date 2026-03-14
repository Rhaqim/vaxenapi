package models

import (
	"time"
)

type BeneficiaryType string

const (
	BeneficiaryTypeBank   BeneficiaryType = "bank"
	BeneficiaryTypeCrypto BeneficiaryType = "crypto"
)

// Beneficiary represents a payment beneficiary
type Beneficiary struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Name           string          `gorm:"not null" json:"name"`
	Type           BeneficiaryType `gorm:"type:text;not null" json:"type"`
	AccountNumber  *string         `json:"accountNumber"`
	RoutingNumber  *string         `json:"routingNumber"`
	BankCode       *string         `json:"bankCode"`
	BankName       *string         `json:"bankName"`
	BankCountry    *string         `gorm:"type:char(2)" json:"bankCountry"`
	Address        *string         `json:"address"`
	Currency       string          `gorm:"type:char(3);not null" json:"currency"`
	Network        *string         `json:"network"`
	IsActive       bool            `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Payouts        []Payout        `gorm:"foreignKey:BeneficiaryID" json:"-"`
}

func (Beneficiary) TableName() string {
	return "beneficiaries"
}
