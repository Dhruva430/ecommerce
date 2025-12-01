package routes

import (
	"api/errors"
	"api/internals/controllers"
	"api/internals/service"
	"api/models/db"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(queries *db.Queries, conn *sql.DB) *gin.Engine {

	r := gin.Default()
	r.Use(errors.GlobalErrorHandler())
	routerAPI := r.Group("/api")

	authService := service.NewAuthService(queries, conn)

	{ // Auth Routes
		authController := controllers.NewAuthController(authService)
		authRoutes := routerAPI.Group("/auth")
		{
			authRoutes.POST("/register", authController.Register)
			authRoutes.POST("/login", authController.Login)
			authRoutes.POST("/refresh-token", authController.RefreshTokenHandler)
		}
	}
	return r
}
