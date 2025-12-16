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

	api := r.Group("/api")

	// Services
	authService := service.NewAuthService(queries, conn)
	userService := service.NewUserService(queries, conn)
	productService := service.NewProductService(queries, conn)
	uploadService := service.NewUploadService(queries, conn)
	orderService := service.NewOrderService(queries, conn)
	sellerService := service.NewSellerService(queries, conn)

	// Controllers
	auth := controllers.NewAuthController(authService)
	users := controllers.NewUserController(userService)
	products := controllers.NewProductController(productService)
	uploads := controllers.NewUploadController(uploadService)
	orders := controllers.NewOrderController(orderService)
	sellers := controllers.NewSellerController(sellerService)

	// ---------------- PUBLIC ROUTES ---------------- //
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", auth.BuyerRegister)
		authRoutes.POST("/login", auth.BuyerLogin)
		authRoutes.POST("/seller/register", sellers.SellerRegisters)
		authRoutes.POST("/seller/login", sellers.LoginSeller)
		authRoutes.POST("/refresh-token", auth.RefreshTokenHandler)
	}

	// ---------------- PROTECTED ROUTES ---------------- //
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	// AUTH
	protected.GET("/me", auth.Me)
	protected.POST("/logout", auth.Logout)

	// USER
	userRoutes := protected.Group("/user")
	{
		userRoutes.DELETE("", users.DeleteUser)
		userRoutes.GET("/addresses", users.GetUserAddress)
		userRoutes.PUT("/addresses", users.UpdateUserAddress)
		userRoutes.GET("/orders", users.GetOrderHistory)
	}

	// PRODUCTS
	productRoutes := protected.Group("/products")
	{
		productRoutes.GET("", products.GetAllProducts)
		productRoutes.GET("/:product_id", products.GetProductByID)
	}

	// SELLER PRODUCT VARIANTS
	protected.POST("/seller/:product_id/variants", products.AddProductVariant)

	// ORDERS
	orderRoutes := protected.Group("/orders")
	{
		orderRoutes.POST("", orders.GetOrder)
		orderRoutes.GET("/:order_id", orders.GetOrderDetails)
		orderRoutes.POST("/:order_id/cancel", orders.CancelOrder)
		orderRoutes.POST("/:order_id/status", orders.UpdateOrderStatus)
	}

	// SELLER
	sellerRoutes := protected.Group("/seller")
	{

		sellerRoutes.POST("/kyc", sellers.ApplyForSellerKYC)

		sellerRoutes.POST("/products", sellers.CreateProduct)
		sellerRoutes.PUT("/products/:product_id", sellers.UpdateProduct)
		sellerRoutes.DELETE("/products/:product_id", sellers.DeleteProduct)
	}

	// UPLOADS
	uploadRoutes := protected.Group("/upload")
	{
		uploadRoutes.POST("/request", uploads.RequestFileUpload)
	}

	return r
}
