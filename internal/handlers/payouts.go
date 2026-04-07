package handlers

import (
	"net/http"

	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// GetPayouts godoc
// @Summary Get all payouts
// @Description Retrieve all payouts for the organization
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /payouts [get]
func (h *Handler) GetPayouts(c *gin.Context) {
	orgID := c.GetString("organizationId")

	payouts, err := h.Services.Payment.GetPayouts(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, payouts)
}

// GetPayout godoc
// @Summary Get payout by ID
// @Description Retrieve a specific payout with latest status from the payment provider
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payout ID"
// @Success 200 {object} map[string]any
// @Router /payouts/{id} [get]
func (h *Handler) GetPayout(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	payout, err := h.Services.Payment.GetPaymentStatus(c.Request.Context(), id, orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, payout)
}

// CreatePayout godoc
// @Summary Send money
// @Description Initiate a payment to a beneficiary. If multi-approval is required, creates an approval request.
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.SendPaymentInput true "Payout details"
// @Success 201 {object} map[string]any
// @Router /payouts [post]
func (h *Handler) CreatePayout(c *gin.Context) {
	var req services.SendPaymentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	req.OrganizationID = c.GetString("organizationId")
	req.InitiatedByID = c.GetString("userId")

	payout, err := h.Services.Payment.InitiatePayment(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, payout)
}

// GetPayoutFees godoc
// @Summary Estimate payout fees
// @Description Get estimated fees for a payment
// @Tags payouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param amount query string true "Amount"
// @Param currency query string true "Currency"
// @Param country query string true "Destination country"
// @Param method query string true "Payment method (wire, ach, sepa, crypto)"
// @Success 200 {object} map[string]any
// @Router /payouts/fees [get]
func (h *Handler) GetPayoutFees(c *gin.Context) {
	amountStr := c.Query("amount")
	currency := c.Query("currency")
	country := c.Query("country")
	method := c.Query("method")

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid amount")
		return
	}

	fees, err := h.Services.Payment.GetFees(c.Request.Context(), amount, currency, country, method)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fees)
}
