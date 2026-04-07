package services_test

import (
	"testing"

	"vaxen/api/internal/models"
	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWebhookService(t *testing.T) *services.WebhookService {
	t.Helper()
	db := testutil.SetupTestDB(t)
	db.AutoMigrate(&models.WebhookEvent{})
	return services.NewWebhookService(db)
}

func TestWebhookService_ProcessWebhook(t *testing.T) {
	svc := setupWebhookService(t)

	event, err := svc.ProcessWebhook("circle", map[string]any{
		"type":           "payment.completed",
		"organizationId": "org-123",
		"amount":         "500.00",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
	assert.Equal(t, "circle", event.Provider)
	assert.Equal(t, "payment.completed", event.EventType)
	assert.Equal(t, models.WebhookStatusPending, event.Status)
	assert.Equal(t, "org-123", event.OrganizationID)
}

func TestWebhookService_ProcessWebhook_UnknownType(t *testing.T) {
	svc := setupWebhookService(t)

	event, err := svc.ProcessWebhook("sumsub", map[string]any{
		"data": "some payload without type",
	})
	require.NoError(t, err)
	assert.Equal(t, "unknown", event.EventType)
}

func TestWebhookService_MarkProcessed(t *testing.T) {
	svc := setupWebhookService(t)

	event, _ := svc.ProcessWebhook("circle", map[string]any{"type": "test"})
	err := svc.MarkProcessed(event.ID)
	require.NoError(t, err)

	events, _ := svc.GetEvents("circle", "processed", 10)
	assert.Len(t, events, 1)
	assert.NotNil(t, events[0].ProcessedAt)
}

func TestWebhookService_MarkFailed(t *testing.T) {
	svc := setupWebhookService(t)

	event, _ := svc.ProcessWebhook("circle", map[string]any{"type": "test"})
	err := svc.MarkFailed(event.ID)
	require.NoError(t, err)

	events, _ := svc.GetEvents("circle", "failed", 10)
	assert.Len(t, events, 1)
	assert.Equal(t, 1, events[0].RetryCount)
}

func TestWebhookService_GetEvents_FilterByProvider(t *testing.T) {
	svc := setupWebhookService(t)

	svc.ProcessWebhook("circle", map[string]any{"type": "a"})
	svc.ProcessWebhook("circle", map[string]any{"type": "b"})
	svc.ProcessWebhook("sumsub", map[string]any{"type": "c"})

	circle, err := svc.GetEvents("circle", "", 50)
	require.NoError(t, err)
	assert.Len(t, circle, 2)

	sumsub, err := svc.GetEvents("sumsub", "", 50)
	require.NoError(t, err)
	assert.Len(t, sumsub, 1)

	all, err := svc.GetEvents("", "", 50)
	require.NoError(t, err)
	assert.Len(t, all, 3)
}

func TestWebhookService_GetEvents_FilterByStatus(t *testing.T) {
	svc := setupWebhookService(t)

	e1, _ := svc.ProcessWebhook("circle", map[string]any{"type": "a"})
	svc.ProcessWebhook("circle", map[string]any{"type": "b"})
	svc.MarkProcessed(e1.ID)

	pending, err := svc.GetEvents("", "pending", 50)
	require.NoError(t, err)
	assert.Len(t, pending, 1)

	processed, err := svc.GetEvents("", "processed", 50)
	require.NoError(t, err)
	assert.Len(t, processed, 1)
}

func TestWebhookService_GetEvents_Limit(t *testing.T) {
	svc := setupWebhookService(t)

	for i := 0; i < 10; i++ {
		svc.ProcessWebhook("circle", map[string]any{"type": "test"})
	}

	events, err := svc.GetEvents("", "", 3)
	require.NoError(t, err)
	assert.Len(t, events, 3)
}
