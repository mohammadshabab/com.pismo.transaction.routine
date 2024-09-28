package utils

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func init() {

}
func AppLog() *slog.Logger {
	envPath := filepath.Join("var/logs", "log.txt")
	logFile, err := os.OpenFile(envPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error while creating log file ", "logPath: ", envPath)
		os.Exit(http.StatusInternalServerError)
	}
	writer := io.Writer(logFile)
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		//	AddSource: true,
	}
	logger := slog.New(
		slog.NewJSONHandler(writer, handlerOpts),
	)
	slog.SetDefault(logger)
	return logger
}
