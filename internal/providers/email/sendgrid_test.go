package email_test

import (
	"context"
	"testing"

	"vaxen/api/internal/providers/email"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendGrid_Name(t *testing.T) {
	sg := email.NewSendGrid(email.SendGridConfig{APIKey: "test"})
	assert.Equal(t, "sendgrid", sg.Name())
}

func TestSendGrid_ImplementsProvider(t *testing.T) {
	var _ email.Provider = email.NewSendGrid(email.SendGridConfig{})
}

func TestSendGrid_Send(t *testing.T) {
	sg := email.NewSendGrid(email.SendGridConfig{APIKey: "test"})
	err := sg.Send(context.Background(), email.Message{
		To:      "user@example.com",
		Subject: "Test",
		HTMLBody: "<p>Hello</p>",
	})
	require.NoError(t, err)
}

func TestSendGrid_SendBatch(t *testing.T) {
	sg := email.NewSendGrid(email.SendGridConfig{APIKey: "test"})
	err := sg.SendBatch(context.Background(), []email.Message{
		{To: "a@example.com", Subject: "A"},
		{To: "b@example.com", Subject: "B"},
	})
	require.NoError(t, err)
}

func TestSendGrid_Defaults(t *testing.T) {
	sg := email.NewSendGrid(email.SendGridConfig{})
	assert.Equal(t, "sendgrid", sg.Name())
}

func TestSMTP_Name(t *testing.T) {
	s := email.NewSMTP(email.SMTPConfig{})
	assert.Equal(t, "smtp", s.Name())
}

func TestSMTP_ImplementsProvider(t *testing.T) {
	var _ email.Provider = email.NewSMTP(email.SMTPConfig{})
}
