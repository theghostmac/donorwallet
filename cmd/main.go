package main

import (
	"github.com/theghostmac/donorwallet/internal/apis"
	"github.com/theghostmac/donorwallet/internal/database"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func main() {
	database.InitDB()
	defer database.CloseDB()

	logger.Info("Starting server on default port 8080...")
	router := apis.InitRouter()
	router.Run(":6569")
}