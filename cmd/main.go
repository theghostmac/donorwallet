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

	port := ":6569"

	logger.Info("Starting server on default port.", zap.String("port", port))
	router := apis.InitRouter()
	router.Run(port)
}