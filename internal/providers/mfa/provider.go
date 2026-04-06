package mfa

import "context"

// Provider defines the interface for multi-factor authentication services.
// Implementations can wrap Twilio Verify, Authy, Google Authenticator (TOTP), etc.
type Provider interface {
	Name() string

	// Enroll starts MFA enrollment for a user. Returns setup data (e.g. TOTP secret, QR URL).
	Enroll(ctx context.Context, req EnrollRequest) (*EnrollResult, error)

	// ConfirmEnrollment verifies the first code to confirm enrollment.
	ConfirmEnrollment(ctx context.Context, userID string, code string) error

	// SendChallenge sends an MFA challenge (SMS, push, etc.). For TOTP this is a no-op.
	SendChallenge(ctx context.Context, req ChallengeRequest) (*ChallengeResult, error)

	// VerifyCode verifies an MFA code provided by the user.
	VerifyCode(ctx context.Context, req VerifyRequest) (*VerifyResult, error)

	// Unenroll removes MFA for a user.
	Unenroll(ctx context.Context, userID string) error
}

type EnrollRequest struct {
	UserID string
	Email  string
	Phone  string // optional, for SMS-based MFA
}

type EnrollResult struct {
	Secret     string // TOTP secret (for authenticator apps)
	QRCodeURL  string // QR code image URL for TOTP
	BackupCode string // One-time backup code
	ProviderID string // Provider's enrollment ID
}

type ChallengeRequest struct {
	UserID  string
	Channel string // "sms", "email", "totp", "push"
}

type ChallengeResult struct {
	ChallengeID string
	ExpiresInS  int // seconds until expiration
}

type VerifyRequest struct {
	UserID      string
	Code        string
	ChallengeID string // optional, depends on provider
}

type VerifyResult struct {
	Valid bool
}
