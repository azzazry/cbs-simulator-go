package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	DatabasePath string
	JWTSecret    string
	Environment  string
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file if exists
	godotenv.Load()

	AppConfig = &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "./database/cbs.db"),
		JWTSecret:    getEnv("JWT_SECRET", "cbs-simulator-secret-key-change-in-production"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}

	log.Printf("Configuration loaded: Port=%s, DB=%s, Env=%s", 
		AppConfig.ServerPort, AppConfig.DatabasePath, AppConfig.Environment)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
