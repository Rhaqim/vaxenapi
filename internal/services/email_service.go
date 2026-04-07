package services

import (
	"context"
	"fmt"

	"vaxen/api/internal/config"
	"vaxen/api/internal/providers"
	"vaxen/api/internal/providers/email"
)

// EmailService sends application emails via the configured email provider.
// Other services call these methods instead of sending emails directly,
// so email content and templates are managed in one place.
type EmailService struct {
	reg         *providers.Registry
	frontendURL string
}

func NewEmailService(cfg *config.Config, reg *providers.Registry) *EmailService {
	return &EmailService{
		reg:         reg,
		frontendURL: cfg.FrontendURL,
	}
}

// SendInvite sends a registration invite to an approved access request applicant.
func (s *EmailService) SendInvite(ctx context.Context, toEmail, toName, inviteToken string) error {
	provider, err := s.reg.RequireEmail()
	if err != nil {
		return err
	}

	registerURL := fmt.Sprintf("%s/register?token=%s", s.frontendURL, inviteToken)

	return provider.Send(ctx, email.Message{
		To:      toEmail,
		ToName:  toName,
		Subject: "You're invited to join Vaxen",
		HTMLBody: fmt.Sprintf(
			`<h2>Welcome to Vaxen</h2>
<p>Hi %s,</p>
<p>Your application has been approved. Click the link below to set up your account:</p>
<p><a href="%s">Complete Registration</a></p>
<p>This link is single-use. If you have any issues, contact support.</p>
<p>— The Vaxen Team</p>`,
			toName, registerURL,
		),
		TextBody: fmt.Sprintf(
			"Hi %s,\n\nYour application has been approved. Visit the link below to set up your account:\n\n%s\n\nThis link is single-use.\n\n— The Vaxen Team",
			toName, registerURL,
		),
		Tags: []string{"invite", "onboarding"},
	})
}

// SendMFAEnrolled notifies a user that MFA has been enabled on their account.
func (s *EmailService) SendMFAEnrolled(ctx context.Context, toEmail, toName string) error {
	provider, err := s.reg.RequireEmail()
	if err != nil {
		return err
	}

	return provider.Send(ctx, email.Message{
		To:      toEmail,
		ToName:  toName,
		Subject: "MFA enabled on your Vaxen account",
		HTMLBody: fmt.Sprintf(
			`<p>Hi %s,</p>
<p>Multi-factor authentication has been enabled on your account. If you did not do this, contact support immediately.</p>
<p>— The Vaxen Team</p>`, toName,
		),
		TextBody: fmt.Sprintf("Hi %s,\n\nMFA has been enabled on your account. If you did not do this, contact support immediately.\n\n— The Vaxen Team", toName),
		Tags:     []string{"security", "mfa"},
	})
}

// SendPayoutInitiated notifies a director that a payout requires approval.
func (s *EmailService) SendPayoutInitiated(ctx context.Context, toEmail, toName, payoutID, amount, currency string) error {
	provider, err := s.reg.RequireEmail()
	if err != nil {
		return err
	}

	approvalURL := fmt.Sprintf("%s/approvals", s.frontendURL)

	return provider.Send(ctx, email.Message{
		To:      toEmail,
		ToName:  toName,
		Subject: fmt.Sprintf("Payout of %s %s requires your approval", amount, currency),
		HTMLBody: fmt.Sprintf(
			`<p>Hi %s,</p>
<p>A payout of <strong>%s %s</strong> has been initiated and requires your approval.</p>
<p><a href="%s">Review and Approve</a></p>
<p>— The Vaxen Team</p>`, toName, amount, currency, approvalURL,
		),
		TextBody: fmt.Sprintf("Hi %s,\n\nA payout of %s %s requires your approval.\n\nReview it at: %s\n\n— The Vaxen Team", toName, amount, currency, approvalURL),
		Tags:     []string{"payout", "approval"},
	})
}

// SendAccessRequestRejected notifies an applicant that their request was rejected.
func (s *EmailService) SendAccessRequestRejected(ctx context.Context, toEmail, toName, reason string) error {
	provider, err := s.reg.RequireEmail()
	if err != nil {
		return err
	}

	reasonText := ""
	if reason != "" {
		reasonText = fmt.Sprintf("<p>Reason: %s</p>", reason)
	}

	return provider.Send(ctx, email.Message{
		To:      toEmail,
		ToName:  toName,
		Subject: "Update on your Vaxen application",
		HTMLBody: fmt.Sprintf(
			`<p>Hi %s,</p>
<p>Thank you for your interest in Vaxen. After reviewing your application, we are unable to proceed at this time.</p>
%s
<p>If you have questions, please contact support.</p>
<p>— The Vaxen Team</p>`, toName, reasonText,
		),
		TextBody: fmt.Sprintf("Hi %s,\n\nThank you for your interest in Vaxen. After reviewing your application, we are unable to proceed at this time.\n\n— The Vaxen Team", toName),
		Tags:     []string{"access-request", "rejection"},
	})
}

// Send is a generic method for sending arbitrary emails via the provider.
func (s *EmailService) Send(ctx context.Context, msg email.Message) error {
	provider, err := s.reg.RequireEmail()
	if err != nil {
		return err
	}
	return provider.Send(ctx, msg)
}
