package middleware

import (
	"ecommerce/utils"
	"log"
	"net/http"
	"strings"

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

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Println("Received Authorization Header:", authHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization token required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if it exists, otherwise use the raw token
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
		// Remove any quotes if present
		tokenString = strings.Trim(tokenString, "\"")

		log.Println("✅ Extracted Token:", tokenString)

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Println("❌ Invalid Token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
