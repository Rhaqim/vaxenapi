package models

import (
	"time"
)

type AccountType string

const (
	AccountTypeIBAN  AccountType = "iban"
	AccountTypePIX   AccountType = "pix"
	AccountTypeACH   AccountType = "ach"
	AccountTypeSWIFT AccountType = "swift"
)

// AccountNumber represents a bank account number
type AccountNumber struct {
	ID             string       `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string       `gorm:"type:uuid;not null" json:"organizationId"`
	Name           string       `gorm:"not null" json:"name"`
	Currency       string       `gorm:"type:char(3);not null" json:"currency"`
	Type           AccountType  `gorm:"type:text;not null" json:"type"`
	AccountNumber  string       `gorm:"not null" json:"accountNumber"`
	RoutingNumber  *string      `json:"routingNumber"`
	BankCode       *string      `json:"bankCode"`
	BankName       string       `gorm:"not null" json:"bankName"`
	BankCountry    string       `gorm:"type:char(2);not null" json:"bankCountry"`
	IsActive       bool         `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (AccountNumber) TableName() string {
	return "account_numbers"
}
