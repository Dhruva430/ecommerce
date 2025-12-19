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

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}
func (u *UserController) GetUserProfile(c *gin.Context) {
	// Implementation for getting user profile
}

func (u *UserController) CreateUserAddress(c *gin.Context) {
	// Implementation for updating user profile
}

func (u *UserController) DeleteUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "unauthorized", Code: http.StatusUnauthorized})
		return
	}
	err := u.service.DeleteUser(c, userID)
	if err != nil {
		c.Error(&errors.AppError{Message: "failed to delete user", Code: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
func (u *UserController) GetUserAddress(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "unauthorized", Code: http.StatusUnauthorized})
		return
	}
	addresses, err := u.service.GetUserAddress(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user address"})
		return
	}
	c.JSON(200, addresses)
}

func (u *UserController) UpdateUserAddress(c *gin.Context) {
	params := c.Params.ByName("address_id")
	AddressID, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid address ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "unauthorized", Code: http.StatusUnauthorized})
		return
	}
	var req request.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := u.service.UpdateUserAddress(c, req, userID, AddressID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user address"})
		return
	}
	c.JSON(200, gin.H{"message": "User address updated successfully"})
}

func (u *UserController) GetOrderHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "unauthorized", Code: http.StatusUnauthorized})
		return
	}

	filter := c.Query("filter")

	orders, err := u.service.GetOrderHistory(c, userID, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}
