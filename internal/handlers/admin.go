package handlers

import (
	"net/http"
	"strconv"

	"vaxen/api/internal/utils"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// GetAuditLogs godoc
// @Summary Get audit logs
// @Description Retrieve audit logs for the organization
// @Tags audit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/audit/logs [get]
func GetAuditLogs(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "log-1",
			"organizationId": organizationID,
			"action":         "user.login",
			"timestamp":      "2026-02-20T10:00:00Z",
			"userId":         "user-123",
		},
	})
}

// --- Platform Admin ---

// GetAllUsers godoc
// @Summary Get all users (Admin)
// @Description Retrieve all users in the system (paginated)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.Services.Admin.GetAllUsers(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(c, users, total, page, limit)
}

// GetUser godoc
// @Summary Get user by ID (Admin)
// @Description Retrieve a specific user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":    id,
		"email": "user@example.com",
		"role":  "user",
	})
}

// UpdateUser godoc
// @Summary Update user (Admin)
// @Description Update a user's information
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "User updated",
	})
}

// DeleteUser godoc
// @Summary Delete user (Admin)
// @Description Delete a user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "User deleted",
	})
}

// GetAllOrganizations godoc
// @Summary Get all organizations (Admin)
// @Description Retrieve all organizations in the system (paginated)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/organizations [get]
func (h *Handler) GetAllOrganizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	orgs, total, err := h.Services.Admin.GetAllOrganizations(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(c, orgs, total, page, limit)
}

// ApproveOrganization godoc
// @Summary Approve organization (Admin)
// @Description Approve a pending organization
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/organizations/{id}/approve [post]
func (h *Handler) ApproveOrganization(c *gin.Context) {
	id := c.Param("id")

	if err := h.Services.Admin.ApproveOrganization(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"status":  "approved",
		"message": "Organization approved",
	})
}

// RejectOrganization godoc
// @Summary Reject organization (Admin)
// @Description Reject a pending organization
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/organizations/{id}/reject [post]
func (h *Handler) RejectOrganization(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.Services.Admin.RejectOrganization(id, req.Reason); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"status":  "rejected",
		"message": "Organization rejected",
	})
}

// --- Platform Settings ---

// GetPlatformSettings godoc
// @Summary Get platform settings (Admin)
// @Description Retrieve all platform settings
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "Filter by category"
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/settings [get]
func (h *Handler) GetPlatformSettings(c *gin.Context) {
	category := c.Query("category")

	settings, err := h.Services.Admin.GetSettings(category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, settings)
}

// UpdatePlatformSetting godoc
// @Summary Update a platform setting (Admin)
// @Description Create or update a platform setting
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/settings [put]
func (h *Handler) UpdatePlatformSetting(c *gin.Context) {
	userID := c.GetString("userId")

	var req struct {
		Key      string `json:"key" binding:"required"`
		Value    string `json:"value" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Services.Admin.UpsertSetting(req.Key, req.Value, req.Category, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Setting updated"})
}

// --- Feature Flags ---

// GetFeatureFlags godoc
// @Summary Get feature flags (Admin)
// @Description Retrieve all feature flags
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/features [get]
func (h *Handler) GetFeatureFlags(c *gin.Context) {
	flags, err := h.Services.Admin.GetFeatureFlags()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, flags)
}

// SetFeatureFlag godoc
// @Summary Set a feature flag (Admin)
// @Description Enable or disable a feature flag
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/features [put]
func (h *Handler) SetFeatureFlag(c *gin.Context) {
	userID := c.GetString("userId")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Enabled     bool   `json:"enabled"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Services.Admin.SetFeatureFlag(req.Name, req.Enabled, req.Description, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Feature flag updated"})
}

// --- Exchange Rates ---

// GetExchangeRates godoc
// @Summary Get exchange rates (Admin)
// @Description Retrieve all active admin-managed exchange rates
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/exchange-rates [get]
func (h *Handler) GetExchangeRates(c *gin.Context) {
	rates, err := h.Services.Admin.GetExchangeRates()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, rates)
}

// UpsertExchangeRate godoc
// @Summary Set an exchange rate (Admin)
// @Description Create or update an exchange rate for the internal exchange provider
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/admin/exchange-rates [put]
func (h *Handler) UpsertExchangeRate(c *gin.Context) {
	userID := c.GetString("userId")

	var req struct {
		FromCurrency string `json:"fromCurrency" binding:"required"`
		ToCurrency   string `json:"toCurrency" binding:"required"`
		Rate         string `json:"rate" binding:"required"`
		Spread       string `json:"spread"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rate, err := decimal.NewFromString(req.Rate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid rate")
		return
	}

	spread := decimal.Zero
	if req.Spread != "" {
		spread, _ = decimal.NewFromString(req.Spread)
	}

	if err := h.Services.Admin.UpsertExchangeRate(req.FromCurrency, req.ToCurrency, rate, spread, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Exchange rate updated"})
}
