package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// GetConversions godoc
// @Summary Get all conversions
// @Description Retrieve all currency conversions for the organization
// @Tags conversions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/conversions [get]
func GetConversions(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "conversion-1",
			"organizationId": organizationID,
			"fromCurrency":   "USD",
			"toCurrency":     "EUR",
			"amount":         "1000.00",
		},
	})
}

// GetConversion godoc
// @Summary Get conversion by ID
// @Description Retrieve a specific conversion
// @Tags conversions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Conversion ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/conversions/{id} [get]
func GetConversion(c *gin.Context) {
	id := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":           id,
		"fromCurrency": "USD",
		"toCurrency":   "EUR",
		"amount":       "1000.00",
	})
}

// CreateConversion godoc
// @Summary Execute a currency swap
// @Description Execute a fiat or crypto conversion using the configured exchange provider
// @Tags conversions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/conversions [post]
func (h *Handler) CreateConversion(c *gin.Context) {
	orgID := c.GetString("organizationId")
	var req struct {
		QuoteID string `json:"quoteId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.Services.Exchange.ExecuteSwap(c.Request.Context(), orgID, req.QuoteID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, order)
}

// CreateQuote godoc
// @Summary Get a conversion quote
// @Description Get a locked-in quote for a currency conversion from the exchange provider
// @Tags quotes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/quotes [post]
func (h *Handler) CreateQuote(c *gin.Context) {
	var req struct {
		FromCurrency string `json:"fromCurrency" binding:"required"`
		ToCurrency   string `json:"toCurrency" binding:"required"`
		Amount       string `json:"amount" binding:"required"`
		Side         string `json:"side" binding:"required"` // "buy" or "sell"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid amount")
		return
	}

	quote, err := h.Services.Exchange.GetQuote(c.Request.Context(), req.FromCurrency, req.ToCurrency, amount, req.Side)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, quote)
}

// GetQuote godoc
// @Summary Get quote by ID
// @Description Retrieve a specific quote
// @Tags quotes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Quote ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/quotes/{id} [get]
func GetQuote(c *gin.Context) {
	id := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":           id,
		"fromCurrency": "USD",
		"toCurrency":   "EUR",
		"rate":         "0.85",
	})
}

// GetSupportedPairs godoc
// @Summary List supported currency pairs
// @Description Get all tradable currency pairs from the exchange provider
// @Tags conversions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/conversions/pairs [get]
func (h *Handler) GetSupportedPairs(c *gin.Context) {
	pairs, err := h.Services.Exchange.ListSupportedPairs(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, pairs)
}

// GetRate godoc
// @Summary Get exchange rate
// @Description Get current exchange rate between two currencies
// @Tags conversions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "Source currency"
// @Param to query string true "Target currency"
// @Success 200 {object} map[string]any
// @Router /api/v1/conversions/rate [get]
func (h *Handler) GetRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "from and to query params required")
		return
	}

	rate, err := h.Services.Exchange.GetRate(c.Request.Context(), from, to)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, rate)
}
