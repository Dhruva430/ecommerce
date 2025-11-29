package service

import (
	"api/errors"
	"net/http"
)

type AuthService struct{}

func NewAuthService() AuthService {
	return AuthService{}
}
func (a *AuthService) Register() error {
	return &errors.AppError{Message: "registration failed", Code: http.StatusBadRequest}
}

func (a *AuthService) Login() error {
	return &errors.AppError{Message: "login failed", Code: http.StatusUnauthorized}
}
