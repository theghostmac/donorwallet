package main

import (
	"os"

	"github.com/theghostmac/donorwallet/internal/apis"
	"github.com/theghostmac/donorwallet/internal/database"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func main() {
	database.InitDB()
	defer database.CloseDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server on default port.", zap.String("port", port))
	router := apis.InitRouter()
	if err := router.Run(":" + port); err != nil {
        logger.Fatal("error: %s", zap.Error(err))
	}
}