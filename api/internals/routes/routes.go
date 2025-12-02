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
	userService := service.NewUserService(queries, conn)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

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
		protected.POST("/logout", authController.Logout)
	}
	{

		protected.DELETE("/user", userController.DeleteUser)
		protected.GET("/user/addresses", userController.GetUserAddress)
		protected.PUT("/user/addresses", userController.UpdateUserAddress)
		protected.GET("/user/orders", userController.GetOrderHistory)
	}

	return r
}
