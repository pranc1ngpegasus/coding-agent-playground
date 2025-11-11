package main

import (
	"flag"
	"os"
	"strconv"
)

// Env holds all configuration variables for the application
type Env struct {
	ServerAddr      string
	OTLPEndpoint    string
	LogLevel        string
	ServiceName     string
	ShutdownTimeout int
}

var env Env

func init() {
	// Define flags
	flag.StringVar(&env.ServerAddr, "addr", getEnvOrDefault("SERVER_ADDR", ":8080"), "Server address")
	flag.StringVar(&env.OTLPEndpoint, "otlp-endpoint", getEnvOrDefault("OTLP_ENDPOINT", "localhost:4317"), "OTLP endpoint")
	flag.StringVar(&env.LogLevel, "log-level", getEnvOrDefault("LOG_LEVEL", "info"), "Log level (debug, info, warn, error)")
	flag.StringVar(&env.ServiceName, "service-name", getEnvOrDefault("SERVICE_NAME", "greet-service"), "Service name for tracing")
	flag.IntVar(&env.ShutdownTimeout, "shutdown-timeout", getEnvOrDefaultInt("SHUTDOWN_TIMEOUT", 30), "Shutdown timeout in seconds")

	flag.Parse()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
