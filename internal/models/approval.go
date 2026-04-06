package models

import (
	"time"

	"gorm.io/datatypes"
)

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
	ApprovalStatusExpired  ApprovalStatus = "expired"
)

type ApprovalActionType string

const (
	ApprovalActionPayout     ApprovalActionType = "payout"
	ApprovalActionConversion ApprovalActionType = "conversion"
	ApprovalActionTransfer   ApprovalActionType = "transfer"
	ApprovalActionAddUser    ApprovalActionType = "add_user"
	ApprovalActionRemoveUser ApprovalActionType = "remove_user"
)

// ApprovalPolicy defines how many director approvals are required for an action type.
type ApprovalPolicy struct {
	ID               string             `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID   string             `gorm:"type:uuid;not null" json:"organizationId"`
	ActionType       ApprovalActionType `gorm:"type:text;not null" json:"actionType"`
	RequiredApprovals int               `gorm:"not null;default:1" json:"requiredApprovals"`
	CreatedAt        time.Time          `json:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	Organization     Organization       `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ApprovalPolicy) TableName() string { return "approval_policies" }

// ApprovalRequest represents a pending action that requires director approval.
type ApprovalRequest struct {
	ID              string         `gorm:"type:uuid;primary_key" json:"id"`
	OrganizationID  string         `gorm:"type:uuid;not null" json:"organizationId"`
	ActionType      ApprovalActionType `gorm:"type:text;not null" json:"actionType"`
	ResourceID      string         `gorm:"type:uuid" json:"resourceId"`    // ID of the payout/conversion/etc.
	ResourceType    string         `gorm:"type:text" json:"resourceType"`  // "payout", "conversion", etc.
	RequestedByID   string         `gorm:"type:uuid;not null" json:"requestedById"`
	Status          ApprovalStatus `gorm:"type:text;default:'pending'" json:"status"`
	RequiredCount   int            `gorm:"not null;default:1" json:"requiredCount"`
	CurrentCount    int            `gorm:"not null;default:0" json:"currentCount"`
	Payload         datatypes.JSON `gorm:"type:jsonb" json:"payload"`      // Snapshot of the action data
	ExpiresAt       *time.Time     `json:"expiresAt"`
	CompletedAt     *time.Time     `json:"completedAt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	Organization    Organization   `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	RequestedBy     User           `gorm:"foreignKey:RequestedByID;references:ID" json:"-"`
	Votes           []ApprovalVote `gorm:"foreignKey:ApprovalRequestID" json:"votes,omitempty"`
}

func (ApprovalRequest) TableName() string { return "approval_requests" }

// ApprovalVote represents a single director's vote on an approval request.
type ApprovalVote struct {
	ID                string         `gorm:"type:uuid;primary_key" json:"id"`
	ApprovalRequestID string         `gorm:"type:uuid;not null" json:"approvalRequestId"`
	VoterID           string         `gorm:"type:uuid;not null" json:"voterId"`
	Decision          ApprovalStatus `gorm:"type:text;not null" json:"decision"` // approved or rejected
	Comment           string         `json:"comment"`
	CreatedAt         time.Time      `json:"createdAt"`
	ApprovalRequest   ApprovalRequest `gorm:"foreignKey:ApprovalRequestID;references:ID" json:"-"`
	Voter             User           `gorm:"foreignKey:VoterID;references:ID" json:"-"`
}

func (ApprovalVote) TableName() string { return "approval_votes" }
