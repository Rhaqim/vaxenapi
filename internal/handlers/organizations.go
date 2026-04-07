package handlers

import (
	"net/http"

	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetOrganizations godoc
// @Summary Get the authenticated user's organization
// @Description Retrieve the organization for the current user
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /organizations [get]
func (h *Handler) GetOrganizations(c *gin.Context) {
	orgID := c.GetString("organizationId")

	org, err := h.Services.Organization.Get(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, org)
}

// GetOrganization godoc
// @Summary Get organization by ID
// @Description Retrieve a specific organization (must be your own)
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]any
// @Router /organizations/{id} [get]
func (h *Handler) GetOrganization(c *gin.Context) {
	id := c.Param("id")

	if id != c.GetString("organizationId") {
		utils.ErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	org, err := h.Services.Organization.Get(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, org)
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Organizations are created during registration
// @Tags organizations
// @Produce json
// @Security BearerAuth
// @Success 400 {object} map[string]any
// @Router /organizations [post]
func (h *Handler) CreateOrganization(c *gin.Context) {
	utils.ErrorResponse(c, http.StatusBadRequest, "organizations are created during registration")
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update organization details
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Param request body services.UpdateOrganizationInput true "Fields to update"
// @Success 200 {object} map[string]any
// @Router /organizations/{id} [put]
func (h *Handler) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")

	if id != c.GetString("organizationId") {
		utils.ErrorResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var req services.UpdateOrganizationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	org, err := h.Services.Organization.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, org)
}
