package kyc

import "context"

// Provider defines the interface for KYC/KYB verification services.
// Implementations can wrap SumSub, Onfido, Jumio, etc.
type Provider interface {
	Name() string

	// --- KYB (Know Your Business) ---

	// SubmitKYB sends business verification data to the provider.
	SubmitKYB(ctx context.Context, req KYBSubmission) (*KYBResult, error)

	// GetKYBStatus retrieves the current status of a KYB application.
	GetKYBStatus(ctx context.Context, referenceID string) (*KYBResult, error)

	// --- KYC (Know Your Customer / Director verification) ---

	// SubmitKYC sends individual identity verification data.
	SubmitKYC(ctx context.Context, req KYCSubmission) (*KYCResult, error)

	// GetKYCStatus retrieves the current status of a KYC application.
	GetKYCStatus(ctx context.Context, referenceID string) (*KYCResult, error)

	// GenerateVerificationURL creates a hosted verification link for a director.
	GenerateVerificationURL(ctx context.Context, applicantID string, redirectURL string) (string, error)
}

// KYBSubmission contains all data needed to submit a business for verification.
type KYBSubmission struct {
	OrganizationID     string
	LegalName          string
	RegistrationNumber string
	TaxID              string
	Country            string
	Address            Address
	Directors          []DirectorInfo
	Documents          []Document
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type DirectorInfo struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
}

type Document struct {
	Type     string `json:"type"` // e.g. "incorporation_cert", "proof_of_address", "id_front"
	FileName string `json:"fileName"`
	Content  []byte `json:"-"`
	URL      string `json:"url,omitempty"`
}

// KYBResult represents the provider's response for a KYB submission.
type KYBResult struct {
	ReferenceID    string            // Provider's reference ID
	Status         VerificationStatus
	Reason         string            // Rejection reason, if any
	ProviderData   map[string]any    // Raw provider-specific data
}

// KYCSubmission contains data for individual identity verification.
type KYCSubmission struct {
	UserID      string
	FirstName   string
	LastName    string
	Email       string
	DateOfBirth string
	Nationality string
	Address     Address
	Documents   []Document
}

// KYCResult represents the provider's response for a KYC submission.
type KYCResult struct {
	ReferenceID  string
	Status       VerificationStatus
	Reason       string
	RiskScore    float64
	ProviderData map[string]any
}

type VerificationStatus string

const (
	StatusPending      VerificationStatus = "pending"
	StatusInReview     VerificationStatus = "in_review"
	StatusApproved     VerificationStatus = "approved"
	StatusRejected     VerificationStatus = "rejected"
	StatusRequiresInfo VerificationStatus = "requires_info"
)
