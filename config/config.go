package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort   string
	DatabasePath string
	Environment  string

	// JWT
	JWTSecret       string
	JWTAccessExpiry  int // minutes
	JWTRefreshExpiry int // hours

	// Security
	RateLimitPerMinute     int
	MaxLoginAttempts       int
	LockoutDurationMinutes int

	// PIN Policy
	PINMinLength int
	PINMaxLength int

	// OTP
	OTPExpiry int // minutes
	OTPLength int
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file if exists
	godotenv.Load()

	AppConfig = &Config{
		// Server
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "./database/cbs.db"),
		Environment:  getEnv("ENVIRONMENT", "development"),

		// JWT
		JWTSecret:       getEnv("JWT_SECRET", "cbs-simulator-secret-key-change-in-production"),
		JWTAccessExpiry:  getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 15),
		JWTRefreshExpiry: getEnvInt("JWT_REFRESH_EXPIRY_HOURS", 168), // 7 days

		// Security
		RateLimitPerMinute:     getEnvInt("RATE_LIMIT_PER_MINUTE", 60),
		MaxLoginAttempts:       getEnvInt("MAX_LOGIN_ATTEMPTS", 3),
		LockoutDurationMinutes: getEnvInt("LOCKOUT_DURATION_MINUTES", 30),

		// PIN Policy
		PINMinLength: getEnvInt("PIN_MIN_LENGTH", 6),
		PINMaxLength: getEnvInt("PIN_MAX_LENGTH", 6),

		// OTP
		OTPExpiry: getEnvInt("OTP_EXPIRY_MINUTES", 5),
		OTPLength: getEnvInt("OTP_LENGTH", 6),
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

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
