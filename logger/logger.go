package logger

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(level string) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level: func(level string) slog.Leveler {
					switch strings.ToUpper(level) {
					case "DEBUG":
						return slog.LevelDebug
					case "INFO":
						return slog.LevelInfo
					case "WARN":
						return slog.LevelWarn
					case "ERROR":
						return slog.LevelError
					default:
						return slog.LevelInfo
					}
				}(level),
			},
		),
	)
}
