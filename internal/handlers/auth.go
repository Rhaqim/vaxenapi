package handlers

import (
	"net/http"

	"vaxen/api/internal/config"
	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// setAuthResponse sets the httpOnly cookie and returns user info.
// The token is never exposed in the response body.
func setAuthResponse(c *gin.Context, cfg *config.Config, result *services.AuthResult) {
	csrfToken, _ := utils.GenerateCSRFToken()
	maxAge := int(cfg.JWTExpiration.Seconds())

	utils.SetAuthCookies(c, result.Token, csrfToken, cfg.SecureCookies, maxAge)

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"csrfToken": csrfToken,
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

// Login godoc
// @Summary User login
// @Description Authenticate user. Token is set as httpOnly cookie, CSRF token returned in body.
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
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		if result.RequiresMFA {
			if req.MFACode == "" {
				challenge, err := svc.MFA.SendChallenge(c.Request.Context(), result.User.ID, "totp")
				if err != nil {
					utils.ErrorResponse(c, http.StatusInternalServerError, "failed to send MFA challenge")
					return
				}
				// Return challenge info without leaking the userId
				utils.SuccessResponse(c, http.StatusOK, gin.H{
					"requiresMfa":  true,
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
				utils.ErrorResponse(c, http.StatusInternalServerError, "authentication failed")
				return
			}
		}

		setAuthResponse(c, cfg, result)
	}
}

// Register godoc
// @Summary Business registration
// @Description Register a new business. Token set as httpOnly cookie.
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

		csrfToken, _ := utils.GenerateCSRFToken()
		maxAge := int(cfg.JWTExpiration.Seconds())
		utils.SetAuthCookies(c, result.Token, csrfToken, cfg.SecureCookies, maxAge)

		utils.SuccessResponse(c, http.StatusCreated, gin.H{
			"csrfToken":   csrfToken,
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

// Logout godoc
// @Summary User logout
// @Description Clear auth cookies
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/v1/auth/logout [post]
func Logout(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.ClearAuthCookies(c, cfg.SecureCookies)
		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"message": "Logged out",
		})
	}
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using the refresh cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/refresh [post]
func RefreshToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement refresh token logic:
		// 1. Read refresh token from httpOnly cookie
		// 2. Validate it against DB (check not revoked)
		// 3. Issue new access token + rotate refresh token
		// 4. Set new cookies
		utils.ErrorResponse(c, http.StatusNotImplemented, "refresh token not yet implemented")
	}
}
