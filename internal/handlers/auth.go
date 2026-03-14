package handlers

import (
	"net/http"

	"vaxen/api/internal/config"
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
	Name         string   `json:"name" binding:"required"`
	Company      string   `json:"company" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	Role         string   `json:"role" binding:"required"`
	Country      string   `json:"country" binding:"required"`
	Markets      []string `json:"markets" binding:"required"`
	AnnualVolume string   `json:"annualVolume" binding:"required"`
	UseCase      string   `json:"useCase" binding:"required"`
	Website      string   `json:"website" binding:"required"`
	Notes        string   `json:"notes"`
	Honeypot     string   `json:"honeypot"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/login [post]
func Login(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: Validate credentials against database
		// For now, returning a mock response

		// Generate JWT token
		token, err := utils.GenerateToken(
			"user-123",
			"org-456",
			req.Email,
			"user",
			cfg.JWTSecret,
			cfg.JWTExpiration,
		)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":             "user-123",
				"email":          req.Email,
				"organizationId": "org-456",
				"role":           "user",
			},
		})
	}
}

// Register godoc
// @Summary User registration
// @Description Register a new user and organization
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /api/v1/auth/register [post]
func Register(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: Create user and organization in database

		// Generate JWT token
		token, err := utils.GenerateToken(
			"new-user-123",
			"new-org-456",
			req.Email,
			"user",
			cfg.JWTSecret,
			cfg.JWTExpiration,
		)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		utils.SuccessResponse(c, http.StatusCreated, gin.H{
			"token": token,
			"user": gin.H{
				"id":             "new-user-123",
				"email":          req.Email,
				"organizationId": "new-org-456",
				"role":           "user",
			},
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
		// TODO: Implement refresh token logic
		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"message": "Token refreshed",
		})
	}
}
