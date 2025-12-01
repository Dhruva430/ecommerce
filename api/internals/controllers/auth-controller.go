package controllers

import (
	"api/errors"
	"api/internals/data/request"
	"api/internals/service"
	"net/http"
	"time"

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
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
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
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errors.AppError{Message: "invalid request", Code: http.StatusBadRequest})
		return
	}
	ip := c.ClientIP()
	req = request.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		IP:       ip,
	}
	user, err := a.service.Login(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Logged in "})

	c.SetCookie("refresh_token", user.RefreshToken, int(time.Until(user.TokenExpiresAt)), "/", "", false, true)
}

func (a *AuthController) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Error(&errors.AppError{Message: "refresh token not found", Code: http.StatusInternalServerError})
		return
	}

	token, err := a.service.RefreshToken(c, refreshToken, c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("refresh_token", token.RefreshToken, int(time.Until(token.TokenExpiresAt)), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed", "access_token": token.AccessToken, "expire_in": token.TokenExpiresAt})
}
