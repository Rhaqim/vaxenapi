package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// ProcessWebhook godoc
// @Summary Process provider webhook
// @Description Receive and store webhook events from payment/KYC providers
// @Tags webhooks
// @Accept json
// @Produce json
// @Param provider path string true "Provider name"
// @Success 200 {object} map[string]any
// @Router /webhooks/provider/{provider} [post]
func (h *Handler) ProcessWebhook(c *gin.Context) {
	providerName := c.Param("provider")
	var payload map[string]any

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Validate webhook signature per provider

	event, err := h.Services.Webhook.ProcessWebhook(providerName, payload)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"eventId":  event.ID,
		"provider": event.Provider,
		"status":   event.Status,
		"message":  "Webhook received",
	})
}

// GetProviders godoc
// @Summary Get active providers
// @Description List all configured provider names
// @Tags providers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /providers [get]
func (h *Handler) GetProviders(c *gin.Context) {
	providerList := []gin.H{}

	if h.Services.KYC != nil {
		providerList = append(providerList, gin.H{"name": "kyc", "status": "configured"})
	}
	if h.Services.MFA != nil {
		providerList = append(providerList, gin.H{"name": "mfa", "status": "configured"})
	}
	if h.Services.Exchange != nil {
		providerList = append(providerList, gin.H{"name": "exchange", "status": "configured"})
	}
	if h.Services.Payment != nil {
		providerList = append(providerList, gin.H{"name": "payment", "status": "configured"})
	}
	if h.Services.Email != nil {
		providerList = append(providerList, gin.H{"name": "email", "status": "configured"})
	}

	utils.SuccessResponse(c, http.StatusOK, providerList)
}

// GetProvider godoc
// @Summary Get provider by name
// @Description Retrieve details about a specific provider
// @Tags providers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Provider name"
// @Success 200 {object} map[string]any
// @Router /providers/{id} [get]
func (h *Handler) GetProvider(c *gin.Context) {
	name := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"name":   name,
		"status": "configured",
	})
}
