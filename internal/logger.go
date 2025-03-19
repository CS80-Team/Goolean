package internal

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	io.Closer
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

func (l *Logger) Warn(prefix, message string, args ...any) {
	message = prefix + message
	l.logger.Warn(message, args...)
}

func (l *Logger) Info(prefix, message string, args ...any) {
	message = prefix + message
	l.logger.Info(message, args...)
}

func (l *Logger) Error(prefix, message string, args ...any) {
	message = prefix + message
	l.logger.Error(message, args...)
}

func (l *Logger) Debug(prefix, message string, args ...any) {
	message = prefix + message
	l.logger.Debug(message, args...)
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}
