package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// ProcessWebhook godoc
// @Summary Process provider webhook
// @Description Process webhook events from payment providers
// @Tags webhooks
// @Accept json
// @Produce json
// @Param provider path string true "Provider name"
// @Success 200 {object} map[string]any
// @Router /api/v1/webhooks/provider/{provider} [post]
func ProcessWebhook(c *gin.Context) {
	provider := c.Param("provider")
	var payload map[string]any

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Validate webhook signature
	// TODO: Process webhook based on provider and event type

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"provider": provider,
		"message":  "Webhook processed successfully",
	})
}

// GetProviders godoc
// @Summary Get all providers
// @Description Retrieve all payment providers
// @Tags providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/providers [get]
func GetProviders(c *gin.Context) {
	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":     "provider-1",
			"name":   "Stripe",
			"type":   "payment",
			"status": "active",
		},
		{
			"id":     "provider-2",
			"name":   "Wise",
			"type":   "transfer",
			"status": "active",
		},
	})
}

// GetProvider godoc
// @Summary Get provider by ID
// @Description Retrieve a specific provider
// @Tags providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Provider ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/providers/{id} [get]
func GetProvider(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":     id,
		"name":   "Stripe",
		"type":   "payment",
		"status": "active",
	})
}
