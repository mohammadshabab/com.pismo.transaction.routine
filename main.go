package main

import (
	"log/slog"

	"com.pismo.transaction.routine/apiroutes"
	"com.pismo.transaction.routine/internal/config"
	"com.pismo.transaction.routine/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	//instantiate DB connection
	database.Connection()
	cfg := config.EnvConfig()
	//Create engine
	route := gin.New()
	apiroutes.AppRoutes(route)
	slog.Info("Starting server on: ", "serverHost: ", cfg.Server.Host, "serverPort: ", cfg.Server.Port)
	route.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}
