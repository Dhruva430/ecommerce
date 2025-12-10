package controllers

import (
	"api/errors"
	"api/internals/data/request"
	"api/internals/middleware"
	"api/internals/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{
		service: service,
	}
}

func (o *OrderController) PlaceOrder(c *gin.Context) {
	var req request.PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
		return
	}
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "user not found in context", Code: http.StatusUnauthorized})
		return
	}
	if err := o.service.PlaceOrder(c, userID, req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order placed successfully"})
}

func (o *OrderController) GetOrderDetails(c *gin.Context) {
	param := c.Param("order_id")
	orderID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.Error(&errors.AppError{Message: "invalid order ID", Code: http.StatusBadRequest})
		return
	}
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "user not found in context", Code: http.StatusUnauthorized})
		return
	}
	order, err := o.service.GetOrderDetails(c, userID, orderID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, order)
}
func (o *OrderController) CancelOrder() {
	// Implementation for canceling an order will go here
}

func (o *OrderController) ListUserOrders() {
	// Implementation for listing user orders will go here
}
func (o *OrderController) UpdateOrderStatus() {
	// Implementation for updating order status will go here
}
func (o *OrderController) ProcessRefund() {
	// Implementation for processing a refund will go here
}
