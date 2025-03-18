package internal

import (
	"log"
	"log/slog"
	"os"
)

type OnClose interface {
	Close() error
}

type Logger struct {
	OnClose
	logger  *slog.Logger
	logFile *os.File
}

func NewLogger(logFilePath string) *Logger {
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	return &Logger{
		logger:  logger,
		logFile: logFile,
	}

}

func (l *Logger) GetLogger() *slog.Logger {
	return l.logger
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}
