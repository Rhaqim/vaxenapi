package models

import (
	"time"

	"gorm.io/datatypes"
)

type AccessRequestStatus string

const (
	AccessRequestPending  AccessRequestStatus = "pending"
	AccessRequestApproved AccessRequestStatus = "approved"
	AccessRequestRejected AccessRequestStatus = "rejected"
)

// AccessRequest stores a business's application to join the platform.
// Submitted via the public request-access form. An admin reviews and
// approves/rejects. On approval, an invite is sent so the director
// can register with a password.
type AccessRequest struct {
	ID           string              `gorm:"type:uuid;primary_key" json:"id"`
	FirstName    string              `gorm:"not null" json:"firstName"`
	LastName     string              `gorm:"not null" json:"lastName"`
	Email        string              `gorm:"not null" json:"email"`
	Company      string              `gorm:"not null" json:"company"`
	Role         string              `gorm:"not null" json:"role"` // applicant's business role
	Country      string              `gorm:"type:char(2)" json:"country"`
	Markets      datatypes.JSON      `gorm:"type:jsonb" json:"markets"`
	AnnualVolume string              `json:"annualVolume"`
	UseCase      string              `json:"useCase"`
	Website      string              `json:"website"`
	Notes        string              `json:"notes"`
	Status       AccessRequestStatus `gorm:"type:text;default:'pending'" json:"status"`
	ReviewedBy   *string             `gorm:"type:uuid" json:"reviewedBy"`
	ReviewedAt   *time.Time          `json:"reviewedAt"`
	ReviewNote   string              `json:"reviewNote"`
	InviteToken  *string             `gorm:"type:text" json:"-"` // set on approval, used for registration
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}

func (AccessRequest) TableName() string { return "access_requests" }
