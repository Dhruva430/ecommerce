package controllers

import (
	"api/internals/data/request"
	"api/internals/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (a *AuthController) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	err := a.service.Register(c, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func (a *AuthController) Login(c *gin.Context) {

	err := a.service.Login()
	if err != nil {
		c.Error(err)
		return
	}
	// c.SetCookie("token", token, util.SESSION_TOKEN_AGE, "/", "", false, true)
	// c.JSON(http.StatusOK, gin.H{"message": "User Logged in as " + user.Username})
}
