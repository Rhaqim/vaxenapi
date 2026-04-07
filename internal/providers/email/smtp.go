package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
)

// SMTP implements the Email Provider interface using plain SMTP.
// Useful for development or self-hosted mail servers.
type SMTP struct {
	host      string
	port      string
	username  string
	password  string
	fromEmail string
	fromName  string
}

type SMTPConfig struct {
	Host      string // e.g. "smtp.gmail.com"
	Port      string // e.g. "587"
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

func NewSMTP(cfg SMTPConfig) *SMTP {
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}
	port := cfg.Port
	if port == "" {
		port = "587"
	}
	fromEmail := cfg.FromEmail
	if fromEmail == "" {
		fromEmail = "noreply@vaxen.io"
	}
	fromName := cfg.FromName
	if fromName == "" {
		fromName = "Vaxen"
	}
	return &SMTP{
		host:      host,
		port:      port,
		username:  cfg.Username,
		password:  cfg.Password,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *SMTP) Name() string { return "smtp" }

func (s *SMTP) Send(ctx context.Context, msg Message) error {
	addr := s.host + ":" + s.port

	body := msg.TextBody
	contentType := "text/plain"
	if msg.HTMLBody != "" {
		body = msg.HTMLBody
		contentType = "text/html"
	}

	headers := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s; charset=UTF-8\r\n\r\n",
		s.fromName, s.fromEmail, msg.To, msg.Subject, contentType,
	)

	fullMsg := []byte(headers + body)

	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	if err := smtp.SendMail(addr, auth, s.fromEmail, []string{msg.To}, fullMsg); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}

	log.Printf("[email:smtp] sent to=%s subject=%q", msg.To, msg.Subject)
	return nil
}

func (s *SMTP) SendBatch(ctx context.Context, msgs []Message) error {
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
