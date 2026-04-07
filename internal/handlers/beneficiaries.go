package handlers

import (
	"net/http"

	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetBeneficiaries godoc
// @Summary Get all beneficiaries
// @Description Retrieve all beneficiaries for the organization
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /beneficiaries [get]
func (h *Handler) GetBeneficiaries(c *gin.Context) {
	orgID := c.GetString("organizationId")

	beneficiaries, err := h.Services.Beneficiary.List(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, beneficiaries)
}

// GetBeneficiary godoc
// @Summary Get beneficiary by ID
// @Description Retrieve a specific beneficiary
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Beneficiary ID"
// @Success 200 {object} map[string]any
// @Router /beneficiaries/{id} [get]
func (h *Handler) GetBeneficiary(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	b, err := h.Services.Beneficiary.Get(orgID, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, b)
}

// CreateBeneficiary godoc
// @Summary Create a new beneficiary
// @Description Add a new beneficiary for the organization
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateBeneficiaryInput true "Beneficiary details"
// @Success 201 {object} map[string]any
// @Router /beneficiaries [post]
func (h *Handler) CreateBeneficiary(c *gin.Context) {
	var req services.CreateBeneficiaryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	req.OrganizationID = c.GetString("organizationId")

	b, err := h.Services.Beneficiary.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, b)
}

// UpdateBeneficiary godoc
// @Summary Update beneficiary
// @Description Update an existing beneficiary
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Beneficiary ID"
// @Param request body services.UpdateBeneficiaryInput true "Fields to update"
// @Success 200 {object} map[string]any
// @Router /beneficiaries/{id} [put]
func (h *Handler) UpdateBeneficiary(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	var req services.UpdateBeneficiaryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	b, err := h.Services.Beneficiary.Update(orgID, id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, b)
}

// DeleteBeneficiary godoc
// @Summary Delete beneficiary
// @Description Soft-delete a beneficiary (marks as inactive)
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Beneficiary ID"
// @Success 200 {object} map[string]any
// @Router /beneficiaries/{id} [delete]
func (h *Handler) DeleteBeneficiary(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	if err := h.Services.Beneficiary.Delete(orgID, id); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Beneficiary deleted",
	})
}
