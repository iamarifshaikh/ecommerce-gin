package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var cartItem models.Cart

	// Bind the incoming JSON
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid Input"})
		return
	}

	// Retrieve the user_id from context (set by AuthRequired)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}
	userID := userIDInterface.(uint)
	cartItem.UserID = userID

	db := database.GetDB()

	// Check if product exists
	var product models.Product
	if err := db.First(&product, cartItem.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Product does not exist"})
		return
	}

	// Check if item is already in cart and update quantity if it is
	var existingCart models.Cart
	if err := db.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).First(&existingCart).Error; err == nil {
		existingCart.Quantity += cartItem.Quantity
		if err := db.Save(&existingCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update cart"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart updated successfully!"})
		return
	}

	// Add new item to cart
	if err := db.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Item added to cart"})
}

func ViewCart(c *gin.Context) {
	var cartItems []models.Cart
	db := database.GetDB()
	userID := c.Param("user_id")

	if err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "cart": cartItems})
}

// Remove item from cart
func RemoveFromCart(c *gin.Context) {
	db := database.GetDB()
	userID := c.Param("user_id")
	cartID := c.Param("cart_id")

	// Ensure user exists
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// Find the cart item and ensure it belongs to the correct user
	var cartItem models.Cart
	if err := db.Where("id = ? AND user_id = ?", cartID, userID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Cart item not found or does not belong to user"})
		return
	}

	// Delete the cart item
	if err := db.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Item removed from cart"})

}
