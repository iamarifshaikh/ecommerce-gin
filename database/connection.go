package database

import (
	"ecommerce/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Construct the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open the database connection
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database! Error: %v", err)
	}

	// Assign to global variable before using it
	DB = database
	fmt.Println("Database connection established!")

	// Auto-migrate database models
	err = DB.AutoMigrate(
		&models.Product{},
		&models.User{},
		&models.Admin{},
		&models.Order{},
		&models.Cart{},
		&models.BlacklistToken{},
	)
	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}
}
