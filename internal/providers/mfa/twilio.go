package mfa

import (
	"context"
	"fmt"
)

// Twilio implements the MFA Provider interface using Twilio Verify.
// See: https://www.twilio.com/docs/verify
type Twilio struct {
	accountSID  string
	authToken   string
	serviceSID  string
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	ServiceSID string // Twilio Verify Service SID
}

func NewTwilio(cfg TwilioConfig) *Twilio {
	return &Twilio{
		accountSID:  cfg.AccountSID,
		authToken:   cfg.AuthToken,
		serviceSID:  cfg.ServiceSID,
	}
}

func (t *Twilio) Name() string { return "twilio" }

func (t *Twilio) Enroll(ctx context.Context, req EnrollRequest) (*EnrollResult, error) {
	// TODO: Implement Twilio Verify - create new factor
	// POST /v2/Services/{ServiceSid}/Entities/{Identity}/NewFactors
	return &EnrollResult{
		ProviderID: fmt.Sprintf("twilio-factor-%s", req.UserID),
		Secret:     "PLACEHOLDER_TOTP_SECRET",
		QRCodeURL:  fmt.Sprintf("otpauth://totp/Vaxen:%s?secret=PLACEHOLDER", req.Email),
	}, nil
}

func (t *Twilio) ConfirmEnrollment(ctx context.Context, userID string, code string) error {
	// TODO: Implement Twilio Verify - verify factor
	// POST /v2/Services/{ServiceSid}/Entities/{Identity}/Factors/{FactorSid}
	return nil
}

func (t *Twilio) SendChallenge(ctx context.Context, req ChallengeRequest) (*ChallengeResult, error) {
	// TODO: Implement Twilio Verify - create challenge
	// POST /v2/Services/{ServiceSid}/Entities/{Identity}/Challenges
	return &ChallengeResult{
		ChallengeID: fmt.Sprintf("twilio-challenge-%s", req.UserID),
		ExpiresInS:  300,
	}, nil
}

func (t *Twilio) VerifyCode(ctx context.Context, req VerifyRequest) (*VerifyResult, error) {
	// TODO: Implement Twilio Verify - verify challenge
	// POST /v2/Services/{ServiceSid}/Entities/{Identity}/Challenges/{ChallengeSid}
	return &VerifyResult{Valid: true}, nil
}

func (t *Twilio) Unenroll(ctx context.Context, userID string) error {
	// TODO: Implement Twilio Verify - delete factor
	// DELETE /v2/Services/{ServiceSid}/Entities/{Identity}/Factors/{FactorSid}
	return nil
}
