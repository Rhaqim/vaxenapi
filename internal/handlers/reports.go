package handlers

import (
	"net/http"
	"strconv"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetStatements godoc
// @Summary Get all statements
// @Description Retrieve all statements for the organization
// @Tags statements
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /statements [get]
func (h *Handler) GetStatements(c *gin.Context) {
	orgID := c.GetString("organizationId")

	statements, err := h.Services.Report.GetStatements(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, statements)
}

// GetStatement godoc
// @Summary Get statement by ID
// @Description Retrieve a specific statement
// @Tags statements
// @Produce json
// @Security BearerAuth
// @Param id path string true "Statement ID"
// @Success 200 {object} map[string]any
// @Router /statements/{id} [get]
func (h *Handler) GetStatement(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	stmt, err := h.Services.Report.GetStatement(orgID, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stmt)
}

// GetFxPnlReport godoc
// @Summary Get FX P&L report
// @Description Retrieve FX profit and loss report for the organization
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /reports/fx-pnl [get]
func (h *Handler) GetFxPnlReport(c *gin.Context) {
	orgID := c.GetString("organizationId")

	report, err := h.Services.Report.GetFxPnlReport(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, report)
}

// GetTransactionReport godoc
// @Summary Get transaction report
// @Description Retrieve transaction report for the organization
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /reports/transactions [get]
func (h *Handler) GetTransactionReport(c *gin.Context) {
	orgID := c.GetString("organizationId")

	report, err := h.Services.Report.GetTransactionReport(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, report)
}

// GetBalanceReport godoc
// @Summary Get balance report
// @Description Retrieve balance report for all wallets
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /reports/balances [get]
func (h *Handler) GetBalanceReport(c *gin.Context) {
	orgID := c.GetString("organizationId")

	entries, err := h.Services.Report.GetBalanceReport(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, entries)
}

// GetAuditLogs godoc
// @Summary Get audit logs
// @Description Retrieve audit logs for the organization (paginated)
// @Tags audit
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} map[string]any
// @Router /audit/logs [get]
func (h *Handler) GetAuditLogs(c *gin.Context) {
	orgID := c.GetString("organizationId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	logs, total, err := h.Services.Report.GetAuditLogs(orgID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(c, logs, total, page, limit)
}
