package controllers

import (
	"bytes"
	"ecommerce/database"
	"ecommerce/models"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProducts retrieves all products
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Products retrieved successfully!",
		"products": products,
	})
}

// GetProductById retrieves a single product
func GetProductById(c *gin.Context) {
	var product models.Product
	if err := database.DB.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Products retrieved successfully!",
		"products": product,
	})
}

// CreateProduct creates a new product
// func CreateProduct(c *gin.Context) {
// 	var product models.Product
// 	if err := c.BindJSON(&product); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
// 		return
// 	}
// 	if err := database.DB.Create(&product).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating product"})
// 		return
// 	}
// 	// Return success response with message
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "Product created successfully!",
// 		"product": product,
// 	})
// }

func CreateProduct(c *gin.Context) {
	var product models.Product

	// Read the raw body for debugging
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("Raw request body: %s\n", string(body))
	// Restore the body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := c.BindJSON(&product); err != nil {
		fmt.Printf("Binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message":       "Invalid request",
			"error":         err.Error(),
			"received_data": string(body),
		})
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully!",
		"product": product,
	})
}

// Update the product
func UpdateProduct(c *gin.Context) {
	// Parse the product ID from URL
	productID := c.Param("id")

	// Check if the product exists
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found!"})
		return
	}

	// Bind the incoming JSON data
	var updatedData models.Product
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input!"})
		return
	}

	// Update the product fields (only the allowed fields)
	product.Name = updatedData.Name
	product.Description = updatedData.Description
	product.Price = updatedData.Price
	product.Stock = updatedData.Stock

	// Save the updated product in the database
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "Product": product})
}

func DeleteProduct(c *gin.Context) {
	// Parese the Product ID from the URL
	productID := c.Param("id")

	// Check if the product exists
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found!"})
		return
	}

	// Delete the product
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully!"})
}
