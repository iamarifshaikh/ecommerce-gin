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

	// User routes

}
