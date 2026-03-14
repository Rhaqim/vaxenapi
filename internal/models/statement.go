package models

import (
	"time"

	"gorm.io/datatypes"
)

type StatementType string

const (
	StatementTypePDF StatementType = "pdf"
	StatementTypeCSV StatementType = "csv"
)

type StatementStatus string

const (
	StatementStatusGenerating StatementStatus = "generating"
	StatementStatusReady      StatementStatus = "ready"
	StatementStatusFailed     StatementStatus = "failed"
)

// StatementFile represents a statement file
type StatementFile struct {
	ID             string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null" json:"organizationId"`
	Type           StatementType   `gorm:"type:text;not null" json:"type"`
	Period         datatypes.JSON  `gorm:"type:jsonb;not null" json:"period"`
	Currency       string          `gorm:"type:char(3);not null" json:"currency"`
	FileURL        string          `gorm:"not null" json:"fileUrl"`
	FileSize       int             `gorm:"not null" json:"fileSize"`
	Status         StatementStatus `gorm:"type:text;default:'generating'" json:"status"`
	GeneratedAt    *time.Time      `json:"generatedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Organization   Organization    `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (StatementFile) TableName() string {
	return "statement_files"
}
