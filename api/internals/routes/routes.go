package routes

import (
	"api/errors"
	"api/internals/controllers"
	"api/internals/service"
	"api/models/db"

	"github.com/gin-gonic/gin"
)

func SetupRouter(queries *db.Queries) *gin.Engine {

	r := gin.Default()
	r.Use(errors.GlobalErrorHandler())
	routerAPI := r.Group("/api")

	authService := service.NewAuthService(queries)

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
