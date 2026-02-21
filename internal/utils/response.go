package utils

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// PaginatedResponse sends a paginated JSON response
func PaginatedResponse(c *gin.Context, data any, total int64, page, limit int) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
