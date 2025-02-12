package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// Init инициализирует логгер
func Init(logLevel, logFormat string) (*slog.Logger, error) {
	var handler slog.Handler

	// Определяем формат логов
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: parseLevel(logLevel),
		})
	case "console":
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: parseLevel(logLevel),
		})
	default:
		return nil, fmt.Errorf("неверный формат log_format в env: %s", logFormat)
	}

	return slog.New(handler), nil
}

// parseLevel преобразует значение из env в константу пакета slog
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
