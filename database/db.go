package database

import (
	"fmt"
	"log"
	"own-paynet/config"
	"own-paynet/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable to hold the database connection
var DB *gorm.DB

// InitDB initializes the database connection and sets up the global DB variable
func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Database connected successfully")
	}

	db.AutoMigrate(&models.User{}, &models.Company{}, &models.Payment{}, &models.PayoutWallet{}, &models.Transaction{}, &models.APIKey{})

	// Set the global DB variable
	DB = db
	return DB
}
