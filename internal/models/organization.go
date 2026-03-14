package models

import (
	"time"

	"gorm.io/datatypes"
)

type KybStatus string

const (
	KybStatusPending      KybStatus = "pending"
	KybStatusApproved     KybStatus = "approved"
	KybStatusRejected     KybStatus = "rejected"
	KybStatusRequiresInfo KybStatus = "requires_info"
)

// Organization represents a business organization
type Organization struct {
	ID                 string              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name               string              `gorm:"not null" json:"name"`
	LegalName          string              `gorm:"not null" json:"legalName"`
	RegistrationNumber string              `gorm:"not null" json:"registrationNumber"`
	TaxID              *string             `json:"taxId"`
	Country            string              `gorm:"type:char(2);not null" json:"country"`
	Address            datatypes.JSON      `gorm:"type:jsonb" json:"address"`
	KYBStatus          KybStatus           `gorm:"type:text;default:'pending'" json:"kybStatus"`
	KYBSubmittedAt     *time.Time          `json:"kybSubmittedAt"`
	KYBApprovedAt      *time.Time          `json:"kybApprovedAt"`
	Settings           datatypes.JSON      `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
	Users              []User              `gorm:"foreignKey:OrganizationID" json:"-"`
	Wallets            []Wallet            `gorm:"foreignKey:OrganizationID" json:"-"`
	AccountNumbers     []AccountNumber     `gorm:"foreignKey:OrganizationID" json:"-"`
	Beneficiaries      []Beneficiary       `gorm:"foreignKey:OrganizationID" json:"-"`
	Deposits           []Deposit           `gorm:"foreignKey:OrganizationID" json:"-"`
	Payouts            []Payout            `gorm:"foreignKey:OrganizationID" json:"-"`
	PayoutBatches      []PayoutBatch       `gorm:"foreignKey:OrganizationID" json:"-"`
	ConversionOrders   []ConversionOrder   `gorm:"foreignKey:OrganizationID" json:"-"`
	LimitOrders        []LimitOrder        `gorm:"foreignKey:OrganizationID" json:"-"`
	Journals           []Journal           `gorm:"foreignKey:OrganizationID" json:"-"`
	ReconciliationRuns []ReconciliationRun `gorm:"foreignKey:OrganizationID" json:"-"`
	StatementFiles     []StatementFile     `gorm:"foreignKey:OrganizationID" json:"-"`
	ComplianceCases    []ComplianceCase    `gorm:"foreignKey:OrganizationID" json:"-"`
	ScreeningResults   []ScreeningResult   `gorm:"foreignKey:OrganizationID" json:"-"`
	AuditLogs          []AuditLog          `gorm:"foreignKey:OrganizationID" json:"-"`
	WebhookEvents      []WebhookEvent      `gorm:"foreignKey:OrganizationID" json:"-"`
}

func (Organization) TableName() string {
	return "organizations"
}
