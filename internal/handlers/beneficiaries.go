package handlers

import (
	"net/http"

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
// @Router /api/v1/beneficiaries [get]
func GetBeneficiaries(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "beneficiary-1",
			"organizationId": organizationID,
			"name":           "John Doe",
			"accountNumber":  "1234567890",
		},
	})
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
// @Router /api/v1/beneficiaries/{id} [get]
func GetBeneficiary(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":            id,
		"name":          "John Doe",
		"accountNumber": "1234567890",
	})
}

// CreateBeneficiary godoc
// @Summary Create a new beneficiary
// @Description Add a new beneficiary
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/beneficiaries [post]
func CreateBeneficiary(c *gin.Context) {
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Create in database
	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-beneficiary-id",
		"message": "Beneficiary created",
	})
}

// UpdateBeneficiary godoc
// @Summary Update beneficiary
// @Description Update an existing beneficiary
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Beneficiary ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/beneficiaries/{id} [put]
func UpdateBeneficiary(c *gin.Context) {
	id := c.Param("id")
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Update in database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Beneficiary updated",
	})
}

// DeleteBeneficiary godoc
// @Summary Delete beneficiary
// @Description Delete a beneficiary
// @Tags beneficiaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Beneficiary ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/beneficiaries/{id} [delete]
func DeleteBeneficiary(c *gin.Context) {
	id := c.Param("id")

	// TODO: Delete from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Beneficiary deleted",
	})
}
