package handlers

import (
	"net/http"

	"vaxen/api/internal/config"
	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token. If MFA is enabled, provide mfaCode.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginInput true "Login credentials"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/login [post]
func Login(cfg *config.Config, svc *services.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.LoginInput
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		result, err := svc.Auth.Login(req)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		if result.RequiresMFA {
			if req.MFACode == "" {
				challenge, err := svc.MFA.SendChallenge(c.Request.Context(), result.User.ID, "totp")
				if err != nil {
					utils.ErrorResponse(c, http.StatusInternalServerError, "failed to send MFA challenge")
					return
				}
				utils.SuccessResponse(c, http.StatusOK, gin.H{
					"requiresMfa":  true,
					"userId":       result.User.ID,
					"challengeId":  challenge.ChallengeID,
					"expiresInSec": challenge.ExpiresInS,
				})
				return
			}

			if err := svc.MFA.VerifyCode(c.Request.Context(), result.User.ID, req.MFACode); err != nil {
				utils.ErrorResponse(c, http.StatusUnauthorized, "invalid MFA code")
				return
			}

			result, err = svc.Auth.CompleteLogin(result.User.ID)
			if err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"token": result.Token,
			"user": gin.H{
				"id":             result.User.ID,
				"email":          result.User.Email,
				"firstName":      result.User.FirstName,
				"lastName":       result.User.LastName,
				"organizationId": result.User.OrganizationID,
				"role":           result.User.Role,
				"isDirector":     result.User.IsDirector,
				"mfaEnabled":     result.User.MFAEnabled,
			},
		})
	}
}

// Register godoc
// @Summary Business registration
// @Description Register a new business organization and its first director
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterInput true "Registration details"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /api/v1/auth/register [post]
func Register(cfg *config.Config, svc *services.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.RegisterInput
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if req.Honeypot != "" {
			utils.SuccessResponse(c, http.StatusCreated, gin.H{"message": "Registration successful"})
			return
		}

		result, err := svc.Auth.Register(req)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		go svc.Wallet.CreateDefaultWallets(c.Request.Context(), result.User.OrganizationID)

		utils.SuccessResponse(c, http.StatusCreated, gin.H{
			"token":       result.Token,
			"requiresMfa": result.RequiresMFA,
			"user": gin.H{
				"id":             result.User.ID,
				"email":          result.User.Email,
				"organizationId": result.User.OrganizationID,
				"role":           result.User.Role,
				"isDirector":     result.User.IsDirector,
			},
			"message": "Registration successful. Please set up MFA to continue.",
		})
	}
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Get a new JWT token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/refresh [post]
func RefreshToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement refresh token logic with token rotation
		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"message": "Token refreshed",
		})
	}
}
