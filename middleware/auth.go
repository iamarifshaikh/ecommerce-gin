package middleware

import (
	"ecommerce/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token not provided!"})
			c.Abort()
			return
		}

		// Check if token is blacklisted
		isBlacklisted, err := utils.IsTokenBlacklisted(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error validating tokens!"})
			c.Abort()
			return
		}

		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or blacklisted token!"})
			c.Abort()
			return
		}

		// Proceed to validate the token normally (existing JWT validation logic)
		c.Next()
	}
}
