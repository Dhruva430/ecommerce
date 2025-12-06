package controllers

import (
	"api/internals/data/request"
	"api/internals/middleware"
	"api/internals/service"

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
