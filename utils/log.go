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
func AppLog() {
	// Ensure the directory exists
	err := os.MkdirAll("/var/logs", os.ModePerm)
	if err != nil {
		slog.Error("Failed to create log directory ", "err: ", err)
		os.Exit(http.StatusInternalServerError)
	}
	envPath := filepath.Join("/var/logs", "log.txt")
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

}
