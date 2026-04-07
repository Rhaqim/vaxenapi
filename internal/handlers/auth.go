package handlers

import (
	"net/http"

	"vaxen/api/internal/config"
	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// setAuthResponse sets the httpOnly cookie and returns user info.
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

// RequestAccess godoc
// @Summary Request platform access
// @Description Submit a business application to join the platform. An admin will review it.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RequestAccessInput true "Access request details"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/request-access [post]
func RequestAccess(svc *services.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.RequestAccessInput
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Honeypot: bots fill hidden fields. Return success so the bot
		// thinks it worked — revealing an error would help it adapt.
		if req.Honeypot != "" {
			utils.SuccessResponse(c, http.StatusCreated, gin.H{
				"message": "Your request has been submitted. We'll be in touch.",
			})
			return
		}

		_, err := svc.Auth.RequestAccess(req)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.SuccessResponse(c, http.StatusCreated, gin.H{
			"message": "Your request has been submitted. We'll be in touch.",
		})
	}
}

// Register godoc
// @Summary Register with invite token
// @Description Create account after an admin has approved your access request. Requires the invite token from the approval email.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterInput true "Registration details"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /auth/register [post]
func Register(cfg *config.Config, svc *services.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.RegisterInput
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
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
// @Router /auth/login [post]
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

// Logout godoc
// @Summary User logout
// @Description Clear auth cookies
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]any
// @Router /auth/logout [post]
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
// @Router /auth/refresh [post]
func RefreshToken(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusNotImplemented, "refresh token not yet implemented")
	}
}
