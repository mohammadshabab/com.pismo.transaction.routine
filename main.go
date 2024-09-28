package main

import (
	"fmt"

	"com.pismo.transaction.routine/apiroutes"
	"com.pismo.transaction.routine/internal/config"
	"com.pismo.transaction.routine/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {

	database.Connection()
	cfg := config.EnvConfig()
	fmt.Printf("Port %s ", cfg.Server.Port)
	route := gin.New()
	apiroutes.AppRoutes(route)
	route.Run(":" + cfg.Server.Port)

	//fmt.Println("db url ", AppConfig.Database.Url)
}
