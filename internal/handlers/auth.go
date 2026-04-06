package handlers

import (
	"net/http"

	"vaxen/api/internal/config"
	"vaxen/api/internal/services"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	MFACode  string `json:"mfaCode,omitempty"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	FirstName          string `json:"firstName" binding:"required"`
	LastName           string `json:"lastName" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	Password           string `json:"password" binding:"required,min=8"`
	CompanyName        string `json:"companyName" binding:"required"`
	LegalName          string `json:"legalName" binding:"required"`
	RegistrationNumber string `json:"registrationNumber" binding:"required"`
	Country            string `json:"country" binding:"required"`
	Honeypot           string `json:"honeypot"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token. If MFA is enabled, provide mfaCode.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/login [post]
func Login(cfg *config.Config, svc *services.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		result, err := svc.Auth.Login(services.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		// MFA required — verify code if provided, otherwise prompt
		if result.RequiresMFA {
			if req.MFACode == "" {
				// Send challenge
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

			// Verify MFA code
			if err := svc.MFA.VerifyCode(c.Request.Context(), result.User.ID, req.MFACode); err != nil {
				utils.ErrorResponse(c, http.StatusUnauthorized, "invalid MFA code")
				return
			}

			// Issue full token
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
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /api/v1/auth/register [post]
func Register(cfg *config.Config, svc *services.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Honeypot check
		if req.Honeypot != "" {
			utils.SuccessResponse(c, http.StatusCreated, gin.H{"message": "Registration successful"})
			return
		}

		result, err := svc.Auth.Register(services.RegisterInput{
			Email:              req.Email,
			Password:           req.Password,
			FirstName:          req.FirstName,
			LastName:           req.LastName,
			CompanyName:        req.CompanyName,
			LegalName:          req.LegalName,
			RegistrationNumber: req.RegistrationNumber,
			Country:            req.Country,
		})
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Create default wallets for the new org
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
