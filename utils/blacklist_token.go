package utils

import (
	"ecommerce/database"
	"ecommerce/models"

	"gorm.io/gorm"
)

func AddTokenToBlacklist(token string) error {
	blacklistToken := models.BlacklistToken{
		Token: token,
	}
	return database.DB.Create(&blacklistToken).Error
}

// IsTokenBlacklisted checks if the token is already blacklisted
func IsTokenBlacklisted(token string) (bool, error) {
	var blacklistToken models.BlacklistToken
	// Query for the first matching token
	err := database.DB.Where("token = ?", token).First(&blacklistToken).Error
	if err != nil {
		// If the token is not found, it's not blacklisted
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		// Handle other errors
		return false, err
	}
	// Token is found in the blacklist
	return true, nil
}
