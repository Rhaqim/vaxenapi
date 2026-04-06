package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.
// In Vaxen, users always belong to an organization. Directors (owners/managers)
// can initiate and approve transactions; other roles have limited permissions.
type User struct {
	ID             string         `gorm:"type:uuid;primary_key" json:"id"`
	Email          string         `gorm:"unique;not null" json:"email"`
	FirstName      string         `gorm:"not null" json:"firstName"`
	LastName       string         `gorm:"not null" json:"lastName"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	Role           UserRole       `gorm:"type:text;default:'viewer'" json:"role"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	IsDirector     bool           `gorm:"default:false" json:"isDirector"`
	MFAEnabled     bool           `gorm:"default:false" json:"mfaEnabled"`
	MFASecret      *string        `json:"-"`
	MFAProvider    string         `gorm:"type:text" json:"-"`
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
