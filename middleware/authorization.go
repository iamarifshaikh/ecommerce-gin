package middleware

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the raw header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("1. Raw Authorization header: %s\n", authHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header not provided"})
			c.Abort()
			return
		}

		// Remove any quotes and trim spaces
		authHeader = strings.Trim(authHeader, "\"")
		authHeader = strings.TrimSpace(authHeader)
		fmt.Printf("2. After removing quotes: %s\n", authHeader)

		// Extract token
		var token string
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		} else {
			token = authHeader
		}
		fmt.Printf("3. Final token for parsing: %s\n", token)

		// Parse token
		claims, err := utils.ParseToken(token)
		if err != nil {
			fmt.Printf("4. Token parsing error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Find user
		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			fmt.Printf("5. Database error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		fmt.Printf("6. Found user: ID=%d, Role=%s\n", user.ID, user.Role)

		if user.Role != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Access restricted to admins"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
