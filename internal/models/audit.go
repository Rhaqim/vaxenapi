package models

import (
	"time"

	"gorm.io/datatypes"
)

type WebhookStatus string

const (
	WebhookStatusPending   WebhookStatus = "pending"
	WebhookStatusProcessed WebhookStatus = "processed"
	WebhookStatusFailed    WebhookStatus = "failed"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	UserID         *string        `gorm:"type:uuid" json:"userId"`
	Action         string         `gorm:"not null" json:"action"`
	Resource       string         `gorm:"not null" json:"resource"`
	ResourceID     *string        `json:"resourceId"`
	Details        datatypes.JSON `gorm:"type:jsonb" json:"details"`
	IPAddress      *string        `json:"ipAddress"`
	UserAgent      *string        `json:"userAgent"`
	CreatedAt      time.Time      `json:"createdAt"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	User           *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	ID             string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrganizationID string         `gorm:"type:uuid;not null" json:"organizationId"`
	Provider       string         `gorm:"not null" json:"provider"`
	EventType      string         `gorm:"not null" json:"eventType"`
	Payload        datatypes.JSON `gorm:"type:jsonb;not null" json:"payload"`
	Status         WebhookStatus  `gorm:"type:text;default:'pending'" json:"status"`
	ProcessedAt    *time.Time     `json:"processedAt"`
	RetryCount     int            `gorm:"default:0" json:"retryCount"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (WebhookEvent) TableName() string {
	return "webhook_events"
}
