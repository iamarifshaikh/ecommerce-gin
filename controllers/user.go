package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers retrieves all users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Users found successfully!",
		"users":   users,
	})
}
