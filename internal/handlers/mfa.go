package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// EnrollMFA godoc
// @Summary Start MFA enrollment
// @Description Begin MFA setup for the authenticated user (required for directors)
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/mfa/enroll [post]
func (h *Handler) EnrollMFA(c *gin.Context) {
	userID := c.GetString("userId")
	email := c.GetString("email")

	result, err := h.Services.MFA.EnrollMFA(c.Request.Context(), userID, email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"qrCodeUrl":  result.QRCodeURL,
		"secret":     result.Secret,
		"backupCode": result.BackupCode,
		"message":    "Scan the QR code with your authenticator app, then confirm with a code.",
	})
}

// ConfirmMFA godoc
// @Summary Confirm MFA enrollment
// @Description Verify first MFA code to complete enrollment
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "MFA code"
// @Success 200 {object} map[string]any
// @Router /api/v1/mfa/confirm [post]
func (h *Handler) ConfirmMFA(c *gin.Context) {
	userID := c.GetString("userId")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Services.MFA.ConfirmMFA(c.Request.Context(), userID, req.Code); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "MFA successfully enabled",
	})
}

// DisableMFA godoc
// @Summary Disable MFA (admin/self)
// @Description Remove MFA enrollment for a user
// @Tags mfa
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /api/v1/mfa/disable [post]
func (h *Handler) DisableMFA(c *gin.Context) {
	userID := c.GetString("userId")

	if err := h.Services.MFA.UnenrollMFA(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "MFA disabled",
	})
}
