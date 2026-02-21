package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetCryptoAddresses godoc
// @Summary Get crypto addresses
// @Description Retrieve all crypto addresses for the organization
// @Tags crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/crypto/addresses [get]
func GetCryptoAddresses(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "address-1",
			"organizationId": organizationID,
			"currency":       "BTC",
			"address":        "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		},
	})
}

// CreateCryptoAddress godoc
// @Summary Create crypto address
// @Description Generate a new crypto address for deposits
// @Tags crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/crypto/addresses [post]
func CreateCryptoAddress(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-address-id",
		"address": "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		"message": "Crypto address created",
	})
}

// CryptoWithdraw godoc
// @Summary Withdraw crypto
// @Description Withdraw cryptocurrency to an external address
// @Tags crypto
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/crypto/withdraw [post]
func CryptoWithdraw(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "withdrawal-123",
		"status":  "pending",
		"message": "Withdrawal initiated",
	})
}
