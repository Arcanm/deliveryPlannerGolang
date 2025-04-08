package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	HTTPPort     int
	GRPCPort     int
	Environment  string
	LogLevel     string
}

func LoadConfig() *Config {
	httpPort, _ := strconv.Atoi(getEnvOrDefault("HTTP_PORT", "8080"))
	grpcPort, _ := strconv.Atoi(getEnvOrDefault("GRPC_PORT", "9090"))

	return &Config{
		MongoURI:     getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnvOrDefault("DB_NAME", "delivery_planner"),
		HTTPPort:     httpPort,
		GRPCPort:     grpcPort,
		Environment:  getEnvOrDefault("ENV", "development"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
