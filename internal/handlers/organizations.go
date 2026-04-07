package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetOrganizations godoc
// @Summary Get all organizations for the authenticated user
// @Description Retrieve all organizations accessible by the current user
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /organizations [get]
func (h *Handler) GetOrganizations(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"organizations": []gin.H{
			{
				"id":     organizationID,
				"name":   "Sample Organization",
				"status": "active",
			},
		},
	})
}

// GetOrganization godoc
// @Summary Get organization by ID
// @Description Retrieve a specific organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]any
// @Router /organizations/{id} [get]
func (h *Handler) GetOrganization(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":     id,
		"name":   "Sample Organization",
		"status": "active",
	})
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Create a new organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /organizations [post]
func (h *Handler) CreateOrganization(c *gin.Context) {
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Create in database
	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-org-id",
		"message": "Organization created",
	})
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update an existing organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]any
// @Router /organizations/{id} [put]
func (h *Handler) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Update in database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Organization updated",
	})
}
