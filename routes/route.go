package routes

import (
	"ecommerce/controllers"
	"ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Register authentication routes
	AuthRoute(router)

	// Product routes
	router.GET("/products", controllers.GetProducts)
	router.GET("/products/:id", controllers.GetProductById)
	router.POST("/product", middleware.AdminOnly(), controllers.CreateProduct)
	router.PUT("/product/:id", middleware.AdminOnly(), controllers.UpdateProduct)
	router.DELETE("/product/:id", middleware.AdminOnly(), controllers.DeleteProduct)

	// Wishlist routes
	router.POST("/wishlist", middleware.AuthRequired(), controllers.AddToWishlist)
	router.GET("/wishlist/:user_id", middleware.AuthRequired(), controllers.ViewWishlist)              // View user's wishlist
	router.DELETE("/wishlist/:wishlist_id", middleware.AuthRequired(), controllers.RemoveFromWishlist) // Remove item from wishlist

	// Cart routes
	router.POST("/cart", controllers.AddToCart)
	router.GET("/cart/:user_id", controllers.ViewCart)
	router.DELETE("/cart/:user_id/:cart_id", controllers.RemoveFromCart)
}
