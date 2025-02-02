package models

import "time"

type Product struct {
	Id          uint    `gorm:"primary_key" json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`

	// Add a foreign key to the Category. This could be a subcategory.
	CategoryID uint      `gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
