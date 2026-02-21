package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// SubmitKYB godoc
// @Summary Submit KYB information
// @Description Submit Know Your Business information for verification
// @Tags kyb
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/kyb/submit [post]
func SubmitKYB(c *gin.Context) {
	organizationID := c.GetString("organizationId")
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"organizationId": organizationID,
		"status":         "pending",
		"message":        "KYB submitted for review",
	})
}

// GetKYBStatus godoc
// @Summary Get KYB status
// @Description Check the status of KYB verification
// @Tags kyb
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/kyb/status [get]
func GetKYBStatus(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"organizationId": organizationID,
		"status":         "approved",
		"verifiedAt":     "2026-01-15T10:00:00Z",
	})
}
