package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email          string         `gorm:"unique;not null" json:"email"`
	FirstName      string         `gorm:"not null" json:"firstName"`
	LastName       string         `gorm:"not null" json:"lastName"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	Role           UserRole       `gorm:"type:text;default:'viewer'" json:"role"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	MFAEnabled     bool           `gorm:"default:false" json:"mfaEnabled"`
	MFASecret      *string        `json:"mfaSecret"`
	LastLoginAt    *time.Time     `json:"lastLoginAt"`
	IsActive       bool           `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	AuditLogs      []AuditLog     `gorm:"foreignKey:UserID" json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
