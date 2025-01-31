package models

import "time"

type Wishlist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
}
