package controllers

import (
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// controllers/authentication.go
func Register(c *gin.Context) {
	// Create a registration request struct
	var registerRequest struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}

	// Bind the incoming JSON to our request struct
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		fmt.Printf("Binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	// Log the received data (remove in production)
	fmt.Printf("Registration Request - Name: %s\n", registerRequest.Name)
	fmt.Printf("Registration Request - Email: %s\n", registerRequest.Email)
	fmt.Printf("Registration Request - Password: %s\n", registerRequest.Password)
	fmt.Printf("Registration Request - Role: %s\n", registerRequest.Role)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error processing password",
		})
		return
	}

	// Create user model
	user := models.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: string(hashedPassword),
		Role:     registerRequest.Role,
	}

	// Validate role
	if user.Role != "" && user.Role != "Admin" && user.Role != "User" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role specified! Please select either Admin or User",
		})
		return
	}

	if user.Role == "" {
		user.Role = "User"
	}

	// Save to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error registering user",
			"details": err.Error(),
		})
		return
	}

	// Test verification immediately
	verifyPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(registerRequest.Password))
	fmt.Printf("Immediate password verification: %v\n", verifyPassword == nil)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind the request body
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials!"})
		return
	}

	// Check for user existence
	var user models.User
	if err := database.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		fmt.Println("User not found in database:", err) // Debug print
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email Doesnt Exist!"})
		return
	}

	// Compare password with hashed password
	fmt.Println("Stored Hashed Password:", user.Password)
	fmt.Println("Entered Password:", credentials.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		fmt.Println("Password comparison failed:", err) // Debug print
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}
	fmt.Println("Password matched successfully!")

	// Generate JWT token
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating JWT token"})
		return
	}

	// Return both tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}
func Logout(c *gin.Context) {
	// Extract the token from the authorization header
	token := c.Request.Header.Get("Authorization")
	fmt.Println("Received Token:", token) // Debugging step

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"token": "Token not provided!"})
		return
	}

	// Check if the token is already blacklisted
	isBlacklisted, err := utils.IsTokenBlacklisted(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check token status!"})
		return
	}

	if isBlacklisted {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is already invalidated!"})
		return
	}

	// Add token to the blacklist
	err = utils.AddTokenToBlacklist(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not logout!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully!"})
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input!"})
		return
	}

	// Validate the refresh token
	claims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired refresh token!"})
		return
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new tokens!"})
		return
	}

	// Send new access and refresh tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":    accessToken,
		"refrheesh_token": newRefreshToken,
	})
}
