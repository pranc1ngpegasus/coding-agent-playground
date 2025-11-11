package main

import (
	loggerPkg "github.com/pranc1ngpegasus/coding-agent-playground/logger"
	"log/slog"
)

type app struct {
	logger *slog.Logger
}

func inject(logLevel string) (*app, error) {
	logger := loggerPkg.NewLogger(logLevel)

	return &app{
		logger: logger,
	}, nil
}
