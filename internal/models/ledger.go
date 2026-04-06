package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type JournalStatus string

const (
	JournalStatusPending  JournalStatus = "pending"
	JournalStatusPosted   JournalStatus = "posted"
	JournalStatusReversed JournalStatus = "reversed"
)

type ReconciliationStatus string

const (
	ReconciliationStatusPending   ReconciliationStatus = "pending"
	ReconciliationStatusCompleted ReconciliationStatus = "completed"
	ReconciliationStatusFailed    ReconciliationStatus = "failed"
)

// Journal represents a journal entry
type Journal struct {
	ID             string        `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string        `gorm:"type:uuid;not null" json:"organizationId"`
	Type           string        `gorm:"not null" json:"type"`
	Description    string        `gorm:"not null" json:"description"`
	Reference      *string       `json:"reference"`
	Status         JournalStatus `gorm:"type:text;default:'pending'" json:"status"`
	PostedAt       *time.Time    `json:"postedAt"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	Organization   Organization  `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Entries        []LedgerEntry `gorm:"foreignKey:JournalID" json:"-"`
}

func (Journal) TableName() string {
	return "journals"
}

// LedgerEntry represents a journal entry line
type LedgerEntry struct {
	ID             string           `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string           `gorm:"type:uuid;not null" json:"organizationId"`
	JournalID      string           `gorm:"type:uuid;not null" json:"journalId"`
	Account        string           `gorm:"not null" json:"account"`
	Debit          *decimal.Decimal `gorm:"type:decimal(20,8)" json:"debit"`
	Credit         *decimal.Decimal `gorm:"type:decimal(20,8)" json:"credit"`
	Description    string           `gorm:"not null" json:"description"`
	Reference      *string          `json:"reference"`
	Metadata       datatypes.JSON   `gorm:"type:jsonb" json:"metadata"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Journal        Journal          `gorm:"foreignKey:JournalID;references:ID" json:"-"`
}

func (LedgerEntry) TableName() string {
	return "ledger_entries"
}

// ReconciliationRun represents a reconciliation run
type ReconciliationRun struct {
	ID             string               `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID string               `gorm:"type:uuid;not null" json:"organizationId"`
	Date           time.Time            `json:"date"`
	Status         ReconciliationStatus `gorm:"type:text;default:'pending'" json:"status"`
	Discrepancies  datatypes.JSON       `gorm:"type:jsonb;default:'[]'" json:"discrepancies"`
	CompletedAt    *time.Time           `json:"completedAt"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	Organization   Organization         `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ReconciliationRun) TableName() string {
	return "reconciliation_runs"
}
