package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetPayouts godoc
// @Summary Get all payouts
// @Description Retrieve all payouts for the organization
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/payouts [get]
func GetPayouts(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "payout-1",
			"organizationId": organizationID,
			"amount":         "500.00",
			"status":         "completed",
		},
	})
}

// GetPayout godoc
// @Summary Get payout by ID
// @Description Retrieve a specific payout
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payout ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/payouts/{id} [get]
func GetPayout(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":     id,
		"amount": "500.00",
		"status": "completed",
	})
}

// CreatePayout godoc
// @Summary Create a new payout
// @Description Create a new payout
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/payouts [post]
func CreatePayout(c *gin.Context) {
	organizationID := c.GetString("organizationId")
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Create in database
	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":             "new-payout-id",
		"organizationId": organizationID,
		"message":        "Payout created",
	})
}
