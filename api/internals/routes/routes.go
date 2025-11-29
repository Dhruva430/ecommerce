package routes

import (
	"api/errors"
	"api/internals/controllers"
	"api/internals/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(errors.GlobalErrorHandler())
	routerAPI := r.Group("/api")

	authService := service.NewAuthService()

	{ // Auth Routes
		authController := controllers.NewAuthController(authService)
		authRoutes := routerAPI.Group("/auth")
		{
			authRoutes.GET("/register", authController.Register)
			authRoutes.GET("/login", authController.Login)
		}
	}
	return r
}
