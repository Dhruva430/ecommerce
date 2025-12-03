package controllers

import (
	"api/internals/data/request"
	"api/internals/middleware"
	"api/internals/service"

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
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return
	}
	var req request.RequestFileUploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := u.service.RequestFileUpload(c, userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to request file upload"})
		return
	}

	c.JSON(200, gin.H{"message": "File upload requested successfully"})

}

func (u *UploadController) DeleteFile() {
	// Implementation for file deletion
}
func (u *UploadController) ListFiles() {
	// Implementation for listing files
}
