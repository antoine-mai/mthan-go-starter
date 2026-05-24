package services

import (
	"log/slog"
	"os"
)

// LoggerService provides structured logging functionality.
type LoggerService struct {
	logger *slog.Logger
}

// NewLoggerService initializes a LoggerService based on the environment string.
// If the environment is "production", it uses a structured JSON handler.
// Otherwise, it defaults to a human-readable text handler.
func NewLoggerService(env string) *LoggerService {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	
	// Also set as default slog logger so any standard slog calls use this handler
	slog.SetDefault(logger)

	return &LoggerService{
		logger: logger,
	}
}

// Info logs a message at LevelInfo.
func (l *LoggerService) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Error logs a message at LevelError.
func (l *LoggerService) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// Debug logs a message at LevelDebug.
func (l *LoggerService) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Warn logs a message at LevelWarn.
func (l *LoggerService) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}
