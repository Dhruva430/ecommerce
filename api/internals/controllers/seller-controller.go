package controllers

import (
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

func (s *SellerController) SellerKYC(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		return
	}

	err := s.service.SellerKYC(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

}
