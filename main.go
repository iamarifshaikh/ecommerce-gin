package main

import (
	"ecommerce/database"
	"ecommerce/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.InitDatabase()

	// Initialize the Gin router
	r := gin.Default()

	// Default route (Hello, Welcome to E-commerce)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, welcome to E-commerce",
		})
	})

	// Setup other routes
	routes.SetupRoutes(r)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}

	log.Println("Server is running on http://localhost:8080")
}
