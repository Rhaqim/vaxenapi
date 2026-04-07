package handlers

import (
	"net/http"

	"vaxen/api/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary Get all orders
// @Description Retrieve all orders for the organization
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "order-1",
			"organizationId": organizationID,
			"status":         "completed",
			"amount":         "100.00",
		},
	})
}

// GetOpenOrders godoc
// @Summary Get open orders
// @Description Retrieve all open orders for the organization
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Router /orders/open [get]
func GetOpenOrders(c *gin.Context) {
	organizationID := c.GetString("organizationId")

	utils.SuccessResponse(c, http.StatusOK, []gin.H{
		{
			"id":             "order-1",
			"organizationId": organizationID,
			"status":         "open",
			"amount":         "100.00",
		},
	})
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Retrieve a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]any
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	id := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":     id,
		"status": "completed",
		"amount": "100.00",
	})
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]any
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"id":      "new-order-id",
		"message": "Order created",
	})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel a pending order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]any
// @Router /orders/{id}/cancel [put]
func CancelOrder(c *gin.Context) {
	id := c.Param("id")

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Order cancelled",
	})
}
