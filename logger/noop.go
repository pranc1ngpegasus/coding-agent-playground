package logger

import "log/slog"

func NewNoopLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
