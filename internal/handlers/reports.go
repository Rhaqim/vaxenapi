package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetStatements godoc
// @Summary Get all statements
// @Description Retrieve all statements for the organization
// @Tags statements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /statements [get]
func GetStatements(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "statement-1",
			"organizationId": organizationID,
			"period":         "2026-01",
			"type":           "account",
		},
	})
}

// GetStatement godoc
// @Summary Get statement by ID
// @Description Retrieve a specific statement
// @Tags statements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Statement ID"
// @Success 200 {object} map[string]any
// @Router /statements/{id} [get]
func GetStatement(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":     id,
		"period": "2026-01",
		"type":   "account",
	})
}

// GetFxPnlReport godoc
// @Summary Get FX P&L report
// @Description Retrieve FX profit and loss report for the organization
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /reports/fx-pnl [get]
func GetFxPnlReport(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Calculate from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"organizationId": organizationID,
		"totalPnl":       "5000.00",
		"currency":       "USD",
		"period":         "2026-01",
	})
}

// GetTransactionReport godoc
// @Summary Get transaction report
// @Description Retrieve transaction report for the organization
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /reports/transactions [get]
func GetTransactionReport(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Calculate from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"organizationId":    organizationID,
		"totalTransactions": 150,
		"totalVolume":       "50000.00",
	})
}

// GetBalanceReport godoc
// @Summary Get balance report
// @Description Retrieve balance report for all accounts
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /reports/balances [get]
func GetBalanceReport(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Calculate from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"organizationId": organizationID,
		"balances": []gin.H{
			{"currency": "USD", "amount": "10000.00"},
			{"currency": "EUR", "amount": "8500.00"},
		},
	})
}
