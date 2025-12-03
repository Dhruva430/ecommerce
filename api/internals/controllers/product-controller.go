package controllers

import (
	"api/errors"
	"api/internals/data/request"
	"api/internals/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

func (p *ProductController) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	sellerID := c.Query("seller_id")
	isActive := c.Query("is_active")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	products, total, err := p.service.GetAllProducts(c, page, pageSize, category, sellerID, isActive)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (p *ProductController) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(&errors.AppError{Message: "invalid product id", Code: http.StatusBadRequest})
		return
	}

	product, err := p.service.GetProductByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (p *ProductController) CreateProduct(c *gin.Context) {
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

	product, err := p.service.CreateProduct(c, req, userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (p *ProductController) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Error(&errors.AppError{Message: "invalid product id", Code: http.StatusBadRequest})
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(&errors.AppError{Message: "user not found in context", Code: http.StatusUnauthorized})
		return
	}

	product, err := p.service.UpdateProduct(c, id, req, userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func (p *ProductController) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
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

	err = p.service.DeleteProduct(c, id, userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
