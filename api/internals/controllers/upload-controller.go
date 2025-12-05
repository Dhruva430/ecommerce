package controllers

import (
	"api/internals/data/request"
	"api/internals/middleware"
	"api/internals/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	service service.UploadService
}

func NewUploadController(service service.UploadService) *UploadController {
	return &UploadController{
		service: service,
	}
}

func (u *UploadController) RequestFileUpload(c *gin.Context) {
	fmt.Print("RequestFileUpload called")
	userID, ok := middleware.GetUserID(c)
	if !ok {
		fmt.Println("Failed to get user ID from context")
		return
	}
	var req request.RequestFileUploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	presignedData, err := u.service.RequestFileUpload(c, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, presignedData)

}
