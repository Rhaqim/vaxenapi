package handlers

import (
	"net/http"

	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetWallets godoc
// @Summary Get all wallets
// @Description Retrieve all fiat wallets for the organization
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /wallets [get]
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
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Success 200 {object} map[string]any
// @Router /wallets/{id} [get]
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
// @Router /wallets [post]
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

	// Fiat wallet
	if req.Currency == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "currency is required for fiat wallets")
		return
	}

	w, err := h.Services.Wallet.CreateFiatWallet(orgID, req.Currency)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, w)
}

// GetWalletBalance godoc
// @Summary Get wallet balance
// @Description Retrieve the current balance of a wallet
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Success 200 {object} map[string]any
// @Router /wallets/{id}/balance [get]
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
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /wallets/web3 [get]
func (h *Handler) GetWeb3Wallets(c *gin.Context) {
	orgID := c.GetString("organizationId")

	wallets, err := h.Services.Wallet.GetWeb3Wallets(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, wallets)
}

// --- Accounts ---

// GetAccounts godoc
// @Summary Get all accounts
// @Description Retrieve all bank accounts for the organization
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /accounts [get]
func (h *Handler) GetAccounts(c *gin.Context) {
	orgID := c.GetString("organizationId")

	accounts, err := h.Services.Account.List(orgID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, accounts)
}

// GetAccount godoc
// @Summary Get account by ID
// @Description Retrieve a specific bank account
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]any
// @Router /accounts/{id} [get]
func (h *Handler) GetAccount(c *gin.Context) {
	orgID := c.GetString("organizationId")
	id := c.Param("id")

	account, err := h.Services.Account.Get(orgID, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, account)
}

// CreateAccount godoc
// @Summary Create a new bank account
// @Description Add a bank account (IBAN, ACH, PIX, SWIFT) for the organization
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateAccountInput true "Account details"
// @Success 201 {object} map[string]any
// @Router /accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var req services.CreateAccountInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	req.OrganizationID = c.GetString("organizationId")

	account, err := h.Services.Account.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, account)
}
