package routes

import (
	"api/errors"
	"api/internals/controllers"
	"api/internals/middleware"
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
	authController := controllers.NewAuthController(authService)

	// -------------------- PUBLIC ROUTES -------------------- //
	authRoutes := routerAPI.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh-token", authController.RefreshTokenHandler)
	}

	// -------------------- PROTECTED ROUTES -------------------- //
	protected := routerAPI.Group("")
	protected.Use(middleware.AuthMiddleware()) // ⬅ access-token validation

	{
		protected.GET("/me", authController.Me)

	}

	return r
}
