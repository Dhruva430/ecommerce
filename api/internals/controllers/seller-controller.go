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

type SellerController struct {
	service service.SellerService
}

func NewSellerController(service service.SellerService) *SellerController {
	return &SellerController{
		service: service,
	}
}

func (s *SellerController) ApplyForSellerKYC(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return
	}
	var req request.SellerKYC
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := s.service.SubmitKYC(c, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

}

func (s *SellerController) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(&errors.AppError{Message: "user not found in context", Code: http.StatusUnauthorized})
		return
	}

	product, err := s.service.CreateProduct(c, req, userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (s *SellerController) UpdateProduct(c *gin.Context) {
	idParam := c.Param("product_id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(&errors.AppError{Message: "invalid product id", Code: http.StatusBadRequest})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.Error(&errors.AppError{Message: "user not found in context", Code: http.StatusUnauthorized})
		return
	}

	product, err := s.service.UpdateProduct(c, productID, req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func (s *SellerController) DeleteProduct(c *gin.Context) {
	idParam := c.Param("product_id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(&errors.AppError{Message: "invalid product id", Code: http.StatusBadRequest})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(&errors.AppError{Message: "user not found in context", Code: http.StatusUnauthorized})
		return
	}

	err = s.service.DeleteProduct(c, id, userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

func (s *SellerController) RegisterSeller(c *gin.Context) {
	var req request.RegisterSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
		return
	}

	err := s.service.RegisterSeller(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seller registered successfully",
	})
}
