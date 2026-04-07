package services

import (
	"encoding/json"
	"errors"
	"time"

	"vaxen/api/internal/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WebhookService struct {
	db *gorm.DB
}

func NewWebhookService(db *gorm.DB) *WebhookService {
	return &WebhookService{db: db}
}

// ProcessWebhook stores a webhook event from a provider and marks it for processing.
func (s *WebhookService) ProcessWebhook(providerName string, payload map[string]any) (*models.WebhookEvent, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("invalid payload")
	}

	eventType, _ := payload["type"].(string)
	if eventType == "" {
		eventType = "unknown"
	}

	// Try to extract org ID from payload (provider-dependent)
	orgID, _ := payload["organizationId"].(string)
	if orgID == "" {
		orgID, _ = payload["organization_id"].(string)
	}

	event := models.WebhookEvent{
		OrganizationID: orgID,
		Provider:       providerName,
		EventType:      eventType,
		Payload:        datatypes.JSON(payloadJSON),
		Status:         models.WebhookStatusPending,
	}
	if err := s.db.Create(&event).Error; err != nil {
		return nil, err
	}

	// TODO: Dispatch to event handlers based on provider + eventType
	// e.g. payment status updates, KYB status changes, etc.

	return &event, nil
}

// MarkProcessed marks a webhook event as processed.
func (s *WebhookService) MarkProcessed(eventID string) error {
	now := time.Now()
	return s.db.Model(&models.WebhookEvent{}).Where("id = ?", eventID).Updates(map[string]any{
		"status":       models.WebhookStatusProcessed,
		"processed_at": now,
	}).Error
}

// MarkFailed marks a webhook event as failed and increments retry count.
func (s *WebhookService) MarkFailed(eventID string) error {
	return s.db.Model(&models.WebhookEvent{}).Where("id = ?", eventID).Updates(map[string]any{
		"status":      models.WebhookStatusFailed,
		"retry_count": gorm.Expr("retry_count + 1"),
	}).Error
}

// GetEvents returns webhook events, optionally filtered by provider and status.
func (s *WebhookService) GetEvents(provider, status string, limit int) ([]models.WebhookEvent, error) {
	q := s.db
	if provider != "" {
		q = q.Where("provider = ?", provider)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if limit <= 0 {
		limit = 50
	}

	var events []models.WebhookEvent
	if err := q.Order("created_at desc").Limit(limit).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
