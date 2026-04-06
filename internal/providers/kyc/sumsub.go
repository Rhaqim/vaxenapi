package kyc

import (
	"context"
	"fmt"
)

// SumSub implements the KYC Provider interface using SumSub's API.
// See: https://docs.sumsub.com/
type SumSub struct {
	apiKey    string
	apiSecret string
	baseURL   string
}

type SumSubConfig struct {
	APIKey    string
	APISecret string
	BaseURL   string // e.g. "https://api.sumsub.com"
}

func NewSumSub(cfg SumSubConfig) *SumSub {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.sumsub.com"
	}
	return &SumSub{
		apiKey:    cfg.APIKey,
		apiSecret: cfg.APISecret,
		baseURL:   baseURL,
	}
}

func (s *SumSub) Name() string { return "sumsub" }

func (s *SumSub) SubmitKYB(ctx context.Context, req KYBSubmission) (*KYBResult, error) {
	// TODO: Implement SumSub API call
	// POST /resources/applicants with level "business-kyb"
	return &KYBResult{
		ReferenceID: fmt.Sprintf("sumsub-kyb-%s", req.OrganizationID),
		Status:      StatusPending,
	}, nil
}

func (s *SumSub) GetKYBStatus(ctx context.Context, referenceID string) (*KYBResult, error) {
	// TODO: Implement SumSub API call
	// GET /resources/applicants/{applicantId}/requiredIdDocsStatus
	return &KYBResult{
		ReferenceID: referenceID,
		Status:      StatusPending,
	}, nil
}

func (s *SumSub) SubmitKYC(ctx context.Context, req KYCSubmission) (*KYCResult, error) {
	// TODO: Implement SumSub API call
	// POST /resources/applicants with level "basic-kyc"
	return &KYCResult{
		ReferenceID: fmt.Sprintf("sumsub-kyc-%s", req.UserID),
		Status:      StatusPending,
	}, nil
}

func (s *SumSub) GetKYCStatus(ctx context.Context, referenceID string) (*KYCResult, error) {
	// TODO: Implement SumSub API call
	return &KYCResult{
		ReferenceID: referenceID,
		Status:      StatusPending,
	}, nil
}

func (s *SumSub) GenerateVerificationURL(ctx context.Context, applicantID string, redirectURL string) (string, error) {
	// TODO: Implement SumSub SDK token generation
	// POST /resources/sdkIntegrations/levels/{levelName}/websdkLink
	return fmt.Sprintf("%s/websdk?applicantId=%s&redirect=%s", s.baseURL, applicantID, redirectURL), nil
}
