package database

import (
	"fmt"
	"log/slog"

	"com.pismo.transaction.routine/internal/config"
	"com.pismo.transaction.routine/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() {
	cfg := config.EnvConfig()
	username := cfg.Database.Username
	password := cfg.Database.Password
	dbHost := cfg.Database.Dbhost
	dbPort := cfg.Database.Dbport
	dbName := cfg.Database.Dbname
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, dbName)
	slog.Info("databse: ", "URI", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, "***", dbHost, dbPort, dbName))
	conn, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
	if err != nil {
		slog.Error("Unable to Establish connection, database error")
	}
	DB = conn
	DB.Debug().AutoMigrate(&models.Account{}, &models.Transaction{})
	slog.Info("connection was successful")
}
