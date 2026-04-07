package email

import (
	"context"
	"fmt"
	"log"
)

// SendGrid implements the Email Provider interface using SendGrid's API.
// See: https://docs.sendgrid.com/api-reference
type SendGrid struct {
	apiKey    string
	fromEmail string
	fromName  string
}

type SendGridConfig struct {
	APIKey    string
	FromEmail string // e.g. "noreply@vaxen.io"
	FromName  string // e.g. "Vaxen"
}

func NewSendGrid(cfg SendGridConfig) *SendGrid {
	fromEmail := cfg.FromEmail
	if fromEmail == "" {
		fromEmail = "noreply@vaxen.io"
	}
	fromName := cfg.FromName
	if fromName == "" {
		fromName = "Vaxen"
	}
	return &SendGrid{
		apiKey:    cfg.APIKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *SendGrid) Name() string { return "sendgrid" }

func (s *SendGrid) Send(ctx context.Context, msg Message) error {
	// TODO: Implement SendGrid v3 Mail Send API
	// POST https://api.sendgrid.com/v3/mail/send
	// With Authorization: Bearer <apiKey>
	log.Printf("[email:sendgrid] sending to=%s subject=%q", msg.To, msg.Subject)
	return nil
}

func (s *SendGrid) SendBatch(ctx context.Context, msgs []Message) error {
	var firstErr error
	for i, msg := range msgs {
		if err := s.Send(ctx, msg); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("failed to send email %d/%d to %s: %w", i+1, len(msgs), msg.To, err)
			}
		}
	}
	return firstErr
}
