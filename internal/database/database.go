package database

import (
	"fmt"
	"log/slog"

	"com.pismo.transaction.routine/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Repo *Repository

type Repository struct {
	RDB *gorm.DB
}

// NewRepository allows for dependency injection of the *gorm.DB connection.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		RDB: db,
	}
}

func Connection(cfg models.AppConfig) {
	//cfg := config.EnvConfig()
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
	Repo = NewRepository(conn)
	// Repo.RDB = conn
	Repo.RDB.Debug().AutoMigrate(&models.Account{}, &models.Transaction{})
	slog.Info("connection was successful")
}
