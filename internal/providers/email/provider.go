package email

import "context"

// Provider defines the interface for sending emails.
// Implementations can wrap SendGrid, AWS SES, Resend, Postmark, etc.
type Provider interface {
	Name() string

	// Send sends a single email.
	Send(ctx context.Context, msg Message) error

	// SendBatch sends multiple emails.
	SendBatch(ctx context.Context, msgs []Message) error
}

// Message represents an email to send.
type Message struct {
	To          string            // Recipient email
	ToName      string            // Recipient display name (optional)
	Subject     string
	TextBody    string            // Plain text body
	HTMLBody    string            // HTML body (preferred if set)
	TemplateID  string            // Provider-hosted template ID (optional)
	TemplateData map[string]string // Template variable substitutions
	ReplyTo     string            // Reply-to address (optional)
	Tags        []string          // Provider tags for analytics (optional)
}
