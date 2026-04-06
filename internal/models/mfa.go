package models

import "time"

// MFAEnrollment tracks a user's MFA enrollment with the active provider.
type MFAEnrollment struct {
	ID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       string     `gorm:"type:uuid;not null;uniqueIndex" json:"userId"`
	Provider     string     `gorm:"type:text;not null" json:"provider"` // e.g. "twilio", "totp"
	ProviderID   string     `gorm:"type:text" json:"providerId"`       // Provider's enrollment/factor ID
	Secret       string     `gorm:"type:text" json:"-"`                // Encrypted TOTP secret
	BackupCode   string     `gorm:"type:text" json:"-"`                // Encrypted backup code
	IsConfirmed  bool       `gorm:"default:false" json:"isConfirmed"`
	ConfirmedAt  *time.Time `json:"confirmedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	User         User       `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (MFAEnrollment) TableName() string { return "mfa_enrollments" }

// MFAChallenge tracks issued MFA challenges for audit and replay protection.
type MFAChallenge struct {
	ID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       string     `gorm:"type:uuid;not null" json:"userId"`
	ChallengeID  string     `gorm:"type:text" json:"challengeId"` // Provider's challenge ID
	Channel      string     `gorm:"type:text;not null" json:"channel"` // "totp", "sms", "email"
	IsVerified   bool       `gorm:"default:false" json:"isVerified"`
	VerifiedAt   *time.Time `json:"verifiedAt"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expiresAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	User         User       `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (MFAChallenge) TableName() string { return "mfa_challenges" }
