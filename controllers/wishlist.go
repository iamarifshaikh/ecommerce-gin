package controllers

import (
	"bytes"
	"ecommerce/database"
	"ecommerce/models"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddToWishlist(c *gin.Context) {
	// Add logging to see the raw request body
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	log.Printf("Raw request body: %s", string(body))

	// Get authenticated user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}

	// Convert userID to uint
	var userIDUint uint
	switch v := userID.(type) {
	case float64:
		userIDUint = uint(v)
	case int:
		userIDUint = uint(v)
	case uint:
		userIDUint = v
	default:
		log.Printf("Unexpected userID type: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Invalid user ID type"})
		return
	}

	// Create a struct to explicitly define the expected request body
	type WishlistRequest struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	var request WishlistRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid Input: product_id is required",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Received request with ProductID: %d", request.ProductID)

	// Create wishlist item
	wishlistItem := models.Wishlist{
		UserID:    userIDUint,
		ProductID: request.ProductID,
	}

	db := database.GetDB()

	// Check if user exists
	var user models.User
	if err := db.First(&user, userIDUint).Error; err != nil {
		log.Printf("User not found with ID: %d", userIDUint)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// Check if product exists
	var product models.Product
	if err := db.First(&product, request.ProductID).Error; err != nil {
		log.Printf("Product not found with ID: %d", request.ProductID)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Product not found"})
		return
	}

	// Prevent duplicate wishlist entries
	var existingWishlist models.Wishlist
	result := db.Where("user_id = ? AND product_id = ?", userIDUint, request.ProductID).First(&existingWishlist)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Item already in wishlist"})
		return
	}

	// Add to wishlist
	if err := db.Create(&wishlistItem).Error; err != nil {
		log.Printf("Error creating wishlist item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Item added to wishlist successfully!",
		"data":    wishlistItem,
	})
}

// View wishlist
func ViewWishlist(c *gin.Context) {
	var wishlist []models.Wishlist
	db := database.GetDB()
	userID := c.Param("user_id")

	// Fetch wishlist with associated product details
	if err := db.Preload("Product").Where("user_id = ?", userID).Find(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "wishlist": wishlist})
}

// Remove item from wishlist
func RemoveFromWishlist(c *gin.Context) {
	var wishlistItem models.Wishlist
	db := database.GetDB()

	if err := db.Where("id = ?", c.Param("wishlist_id")).Delete(&wishlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item removed from wishlist"})
}
