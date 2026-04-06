package middleware

import (
	"net/http"
	"strings"

	"vaxen/api/internal/config"
	"vaxen/api/internal/models"
	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens from httpOnly cookie or Authorization header.
// For cookie-based auth on state-changing methods (POST/PUT/DELETE), it also
// validates the X-CSRF-Token header matches the csrf cookie to prevent CSRF attacks.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, source := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// CSRF check: required for cookie-based auth on mutating requests
		if source == "cookie" && isMutatingMethod(c.Request.Method) {
			csrfHeader := c.GetHeader(utils.CSRFTokenHeader)
			csrfCookie, err := c.Cookie(utils.CSRFTokenCookie)
			if err != nil || csrfHeader == "" || csrfHeader != csrfCookie {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or missing CSRF token"})
				c.Abort()
				return
			}
		}

		c.Set("userId", claims.UserID)
		c.Set("organizationId", claims.OrganizationID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func OptionalAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := extractToken(c)
		if tokenString == "" {
			c.Next()
			return
		}

		if claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret); err == nil {
			c.Set("userId", claims.UserID)
			c.Set("organizationId", claims.OrganizationID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}

// AdminMiddleware ensures the user has the platform admin role.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if models.UserRole(role) != models.UserRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// DirectorMiddleware ensures the user is a director (owner or manager).
func DirectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := models.UserRole(c.GetString("role"))
		if !role.CanApprove() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Director access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// extractToken tries to get the JWT from:
// 1. httpOnly cookie (browser clients)
// 2. Authorization: Bearer header (API clients)
// Returns the token string and the source ("cookie" or "header").
func extractToken(c *gin.Context) (string, string) {
	// Try cookie first
	if token, err := c.Cookie(utils.AccessTokenCookie); err == nil && token != "" {
		return token, "cookie"
	}

	// Fallback to Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ""
	}

	return parts[1], "header"
}

func isMutatingMethod(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE"
}
