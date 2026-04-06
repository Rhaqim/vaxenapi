package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetWallets godoc
// @Summary Get all wallets
// @Description Retrieve all fiat wallets for the organization
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/wallets [get]
func (h *Handler) GetWallets(c *gin.Context) {
	orgID := c.GetString("organizationId")

	wallets, err := h.Services.Wallet.GetWallets(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, wallets)
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
func (h *Handler) GetWallet(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	wallet, err := h.Services.Wallet.GetWallet(orgID, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, wallet)
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new fiat or Web3 wallet for the organization
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /api/v1/wallets [post]
func (h *Handler) CreateWallet(c *gin.Context) {
	orgID := c.GetString("organizationId")
	var req struct {
		Type     string `json:"type" binding:"required"` // "fiat" or "crypto"
		Currency string `json:"currency"`                // For fiat wallets
		Network  string `json:"network"`                 // For crypto wallets
		Label    string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Type == "crypto" {
		w, err := h.Services.Wallet.CreateWeb3Wallet(c.Request.Context(), orgID, req.Network, req.Label)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SuccessResponse(c, http.StatusCreated, w)
		return
	}

	// Fiat wallet creation handled via DB directly for now
	utils.SuccessResponse(c, http.StatusCreated, gin.H{
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
func (h *Handler) GetWalletBalance(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	wallet, err := h.Services.Wallet.GetWallet(orgID, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"walletId":         wallet.ID,
		"balance":          wallet.Balance,
		"availableBalance": wallet.AvailableBalance,
		"pendingBalance":   wallet.PendingBalance,
		"currency":         wallet.Currency,
	})
}

// GetWeb3Wallets godoc
// @Summary Get Web3 wallets
// @Description Retrieve all blockchain wallets for the organization
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/wallets/web3 [get]
func (h *Handler) GetWeb3Wallets(c *gin.Context) {
	orgID := c.GetString("organizationId")

	wallets, err := h.Services.Wallet.GetWeb3Wallets(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, wallets)
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
