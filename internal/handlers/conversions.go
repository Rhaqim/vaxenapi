package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

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
// @Summary Create a new conversion
// @Description Create a new currency conversion
// @Tags conversions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/conversions [post]
func CreateConversion(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-conversion-id",
		"message": "Conversion created",
	})
}

// CreateQuote godoc
// @Summary Create a new quote
// @Description Get a quote for currency conversion
// @Tags quotes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/quotes [post]
func CreateQuote(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":           "quote-123",
		"fromCurrency": "USD",
		"toCurrency":   "EUR",
		"rate":         "0.85",
		"expiresAt":    "2026-02-20T12:00:00Z",
	})
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
