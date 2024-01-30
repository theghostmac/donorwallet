package database

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/theghostmac/donorwallet/internal/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewDevelopment()
	DB *gorm.DB
)

func InitDB() {
    logger.Info("Initializing database...")

    var err error
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        logger.Fatal("DATABASE_URL environment variable not set.")
    }

    DB, err = gorm.Open("postgres", dbURL)
    if err != nil {
        logger.Fatal("Failed to connect to database.", zap.Error(err))
    }

    logger.Info("Successfully connected to the database.")
    DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{}, &models.Donation{})
}


// CloseDB closes the database conection.
func CloseDB() {
	logger.Info("Closing database...")
	err := DB.Close()
	if err!= nil {
		logger.Fatal("Failed to close database.", zap.Error(err))
	}
}