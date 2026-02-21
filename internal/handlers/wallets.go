package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetWallets godoc
// @Summary Get all wallets
// @Description Retrieve all wallets for the organization
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/wallets [get]
func GetWallets(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "wallet-1",
			"organizationId": organizationID,
			"currency":       "USD",
			"balance":        "1000.00",
		},
	})
}

// GetWallet godoc
// @Summary Get wallet by ID
// @Description Retrieve a specific wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/wallets/{id} [get]
func GetWallet(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":       id,
		"currency": "USD",
		"balance":  "1000.00",
	})
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet for the organization
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/wallets [post]
func CreateWallet(c *gin.Context) {
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Create in database
	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-wallet-id",
		"message": "Wallet created",
	})
}

// GetWalletBalance godoc
// @Summary Get wallet balance
// @Description Retrieve the current balance of a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/wallets/{id}/balance [get]
func GetWalletBalance(c *gin.Context) {
	id := c.Param("id")

	// TODO: Calculate from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"walletId": id,
		"balance":  "1000.00",
		"currency": "USD",
	})
}

// GetAccounts godoc
// @Summary Get all accounts
// @Description Retrieve all accounts for the organization
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/accounts [get]
func GetAccounts(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "account-1",
			"organizationId": organizationID,
			"accountNumber":  "1234567890",
			"currency":       "USD",
		},
	})
}

// GetAccount godoc
// @Summary Get account by ID
// @Description Retrieve a specific account
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]any
// @Router /api/v1/accounts/{id} [get]
func GetAccount(c *gin.Context) {
	id := c.Param("id")

	// TODO: Fetch from database
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":            id,
		"accountNumber": "1234567890",
		"currency":      "USD",
	})
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account for the organization
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/accounts [post]
func CreateAccount(c *gin.Context) {
	var req map[string]any

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Create in database
	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-account-id",
		"message": "Account created",
	})
}
