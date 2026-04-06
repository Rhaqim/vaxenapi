package handlers

import (
	"net/http"

	"vaxen/api/internal/providers/kyc"
	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

type KYBSubmitRequest struct {
	LegalName          string             `json:"legalName" binding:"required"`
	RegistrationNumber string             `json:"registrationNumber" binding:"required"`
	TaxID              string             `json:"taxId"`
	Country            string             `json:"country" binding:"required"`
	Address            kyc.Address        `json:"address" binding:"required"`
	Directors          []kyc.DirectorInfo `json:"directors" binding:"required"`
}

// SubmitKYB godoc
// @Summary Submit KYB information
// @Description Submit Know Your Business information for verification via the configured KYC provider
// @Tags kyb
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body KYBSubmitRequest true "KYB data"
// @Success 201 {object} map[string]any
// @Router /api/v1/kyb/submit [post]
func (h *Handler) SubmitKYB(c *gin.Context) {
	orgID := c.GetString("organizationId")
	var req KYBSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.Services.KYC.SubmitKYB(c.Request.Context(), services.KYBSubmitInput{
		OrganizationID:     orgID,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		TaxID:              req.TaxID,
		Country:            req.Country,
		Address:            req.Address,
		Directors:          req.Directors,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"referenceId":    result.ReferenceID,
		"status":         result.Status,
		"organizationId": orgID,
		"message":        "KYB submitted for review",
	})
}

// GetKYBStatus godoc
// @Summary Get KYB status
// @Description Check the status of KYB verification from the provider
// @Tags kyb
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/kyb/status [get]
func (h *Handler) GetKYBStatus(c *gin.Context) {
	orgID := c.GetString("organizationId")

	result, err := h.Services.KYC.GetKYBStatus(c.Request.Context(), orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"organizationId": orgID,
		"status":         result.Status,
		"referenceId":    result.ReferenceID,
		"reason":         result.Reason,
	})
}
