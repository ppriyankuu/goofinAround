package main

import (
	"basic-api/api"
	"basic-api/config"
	"basic-api/middlewares"

	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := api.SetupRouter()

	r.Use(middlewares.CORS())
	r.Use(middlewares.Logging(logger))
	r.Use(middlewares.Recovery(logger))

	// Start server
	r.Run(":" + config.AppPort())
}
