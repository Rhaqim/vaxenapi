package models

import (
	"time"

	"gorm.io/datatypes"
)

type ComplianceType string

const (
	ComplianceTypeKYB ComplianceType = "kyb"
	ComplianceTypeKYT ComplianceType = "kyt"
	ComplianceTypeAML ComplianceType = "aml"
)

type ComplianceStatus string

const (
	ComplianceStatusPending      ComplianceStatus = "pending"
	ComplianceStatusApproved     ComplianceStatus = "approved"
	ComplianceStatusRejected     ComplianceStatus = "rejected"
	ComplianceStatusRequiresInfo ComplianceStatus = "requires_info"
)

type ScreeningType string

const (
	ScreeningTypeTransaction ScreeningType = "transaction"
	ScreeningTypeAddress     ScreeningType = "address"
	ScreeningTypeEntity      ScreeningType = "entity"
)

type ScreeningStatus string

const (
	ScreeningStatusClean   ScreeningStatus = "clean"
	ScreeningStatusFlagged ScreeningStatus = "flagged"
	ScreeningStatusBlocked ScreeningStatus = "blocked"
)

// ComplianceCase represents a compliance case
type ComplianceCase struct {
	ID             string           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string           `gorm:"type:uuid;not null" json:"organizationId"`
	Type           ComplianceType   `gorm:"type:text;not null" json:"type"`
	Status         ComplianceStatus `gorm:"type:text;default:'pending'" json:"status"`
	Provider       string           `gorm:"not null" json:"provider"`
	ProviderCaseID string           `gorm:"not null" json:"providerCaseId"`
	SubmittedAt    time.Time        `json:"submittedAt"`
	CompletedAt    *time.Time       `json:"completedAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Organization   Organization     `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ComplianceCase) TableName() string {
	return "compliance_cases"
}

// ScreeningResult represents a screening result
type ScreeningResult struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Type           ScreeningType   `gorm:"type:text;not null" json:"type"`
	Status         ScreeningStatus `gorm:"type:text;default:'clean'" json:"status"`
	RiskScore      int             `gorm:"not null" json:"riskScore"`
	Provider       string          `gorm:"not null" json:"provider"`
	Details        datatypes.JSON  `gorm:"type:jsonb;not null" json:"details"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ScreeningResult) TableName() string {
	return "screening_results"
}
