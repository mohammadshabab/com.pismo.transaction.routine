package config

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"com.pismo.transaction.routine/models"
	"com.pismo.transaction.routine/utils"
	"github.com/spf13/viper"
)

var configErr slog.Attr

func init() {
	configErr = slog.Group(
		"configErr",
		"file", "config.go",
	)

}
func getConfig(log *slog.Logger) (models.AppConfig, error) {
	conf := viper.New()
	envPath := filepath.Join("environment", "config.yaml")
	conf.SetConfigFile(envPath)

	// Set up environment variable key replacement
	replacer := strings.NewReplacer(".", "_")
	conf.SetEnvKeyReplacer(replacer)
	conf.AutomaticEnv()

	// Read the configuration file
	err := conf.ReadInConfig()
	if err != nil {
		log.Error("Error reading config file", "error", err, configErr)
		return models.AppConfig{}, err // Return empty config and error
	}

	var cfg models.AppConfig

	// Unmarshal the config into AppConfig struct
	if err := conf.Unmarshal(&cfg); err != nil {
		log.Error("Error unmarshalling config", "error", err, configErr)
		return models.AppConfig{}, err // Return empty config and error
	}

	return cfg, nil // Return the config and nil error if successful
}

func EnvConfig() models.AppConfig {
	logger := utils.AppLog()
	AppConfig, err := getConfig(logger)
	if err != nil {
		slog.Error("Unable to fetch config ")
		os.Exit(http.StatusInternalServerError)
	}
	return AppConfig
}
