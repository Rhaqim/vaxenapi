package handlers

import (
	"net/http"

	"vaxen/api/internal/models"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetPendingApprovals godoc
// @Summary Get pending approval requests
// @Description Retrieve all pending approval requests for the organization
// @Tags approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /approvals/pending [get]
func (h *Handler) GetPendingApprovals(c *gin.Context) {
	orgID := c.GetString("organizationId")

	requests, err := h.Services.Approval.GetPendingRequests(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, requests)
}

// GetApprovalRequest godoc
// @Summary Get approval request by ID
// @Description Retrieve a specific approval request with all votes
// @Tags approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Approval Request ID"
// @Success 200 {object} map[string]any
// @Router /approvals/{id} [get]
func (h *Handler) GetApprovalRequest(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	request, err := h.Services.Approval.GetRequest(id, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, request)
}

// VoteOnApproval godoc
// @Summary Vote on an approval request
// @Description Submit a director's approval or rejection vote
// @Tags approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Approval Request ID"
// @Success 200 {object} map[string]any
// @Router /approvals/{id}/vote [post]
func (h *Handler) VoteOnApproval(c *gin.Context) {
	userID := c.GetString("userId")
	id := c.Param("id")

	var req struct {
		Decision string `json:"decision" binding:"required"` // "approved" or "rejected"
		Comment  string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	decision := models.ApprovalStatus(req.Decision)
	if decision != models.ApprovalStatusApproved && decision != models.ApprovalStatusRejected {
		utils.ErrorResponse(c, http.StatusBadRequest, "decision must be 'approved' or 'rejected'")
		return
	}

	if err := h.Services.Approval.Vote(c.Request.Context(), id, userID, decision, req.Comment); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Vote recorded",
	})
}

// UpdateApprovalPolicy godoc
// @Summary Update approval policy
// @Description Set the number of required approvals for an action type (business admin)
// @Tags approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /approvals/policy [put]
func (h *Handler) UpdateApprovalPolicy(c *gin.Context) {
	orgID := c.GetString("organizationId")

	var req struct {
		ActionType        string `json:"actionType" binding:"required"`
		RequiredApprovals int    `json:"requiredApprovals" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Services.Approval.UpdatePolicy(orgID, models.ApprovalActionType(req.ActionType), req.RequiredApprovals); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Approval policy updated",
	})
}
