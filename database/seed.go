package database

import (
	"ecommerce/models"
	"log"
)

func SeedCategories() {
	db := GetDB()

	// Define top-level categories with subcategories (example for "Mobiles, Computers")
	// You can add additional top-level categories and their subcategories similarly.
	categories := []models.Category{
		{
			Name: "Mobiles, Computers",
			Subcategories: []models.Category{
				{Name: "Mobiles, Tablets & More"},
				{Name: "All Mobile Phones"},
				{Name: "All Mobile Accessories"},
				{Name: "Cases & Covers"},
				{Name: "Screen Protectors"},
				{Name: "Power Banks"},
				{Name: "Refurbished & Open Box"},
				{Name: "Tablets"},
				{Name: "Wearable Devices"},
				{Name: "Smart Home"},
				{Name: "Office Supplies & Stationery"},
				{Name: "Software"},
				{Name: "Computers & Accessories"},
				{Name: "All Computers & Accessories"},
				{Name: "Laptops"},
				{Name: "Drives & Storage"},
				{Name: "Printers & Ink"},
				{Name: "Networking Devices"},
				{Name: "Computer Accessories"},
				{Name: "Game Zone"},
				{Name: "Monitors"},
				{Name: "Desktops"},
				{Name: "Components"},
				{Name: "All Electronics"}, // You might want to verify if "All Electronics" is needed here or belongs elsewhere
			},
		},
		{
			Name: "TV, Appliances, Electronics",
		},
		{
			Name: "Men's Fashion",
		},
		{
			Name: "Women's Fashion",
		},
		{
			Name: "Home, Kitchen, Pets",
		},
		{
			Name: "Beauty, Health, Grocery",
		},
		{
			Name: "Sports, Fitness, Bags, Luggage",
		},
		{
			Name: "Toys, Baby Products, Kids' Fashion",
		},
		{
			Name: "Car, Motorbike, Industrial",
		},
		{
			Name: "Books",
		},
		{
			Name: "Movies, Music & Video Games",
		},
	}

	// Iterate over the defined categories and insert them if they don't exist.
	for _, cat := range categories {
		var existing models.Category
		if err := db.Where("name = ?", cat.Name).First(&existing).Error; err != nil {
			if err := db.Create(&cat).Error; err != nil {
				log.Printf("Error seeding category %s: %v", cat.Name, err)
			}
		}
	}
}
