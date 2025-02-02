package database

import (
	"ecommerce/models" // Import your models
	"fmt"
	"log"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// Fetch database credentials from environment variables
	username := os.Getenv("DB_USERNAME")
	password := url.QueryEscape(os.Getenv("DB_PASSWORD")) // URL encode the password
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Check if any required environment variable is missing
	if username == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Missing required environment variables for database connection")
	}

	// Construct the connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Shanghai",
		host, username, password, dbname, port)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations for all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Cart{},
		&models.BlacklistToken{},
		&models.Admin{},
		&models.Wishlist{},
		&models.Category{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Set the global DB variable
	DB = db
	log.Println("Successfully connected and migrated the database!")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
