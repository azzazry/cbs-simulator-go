package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort  string
	DatabaseDSN string
	Environment string

	// JWT
	JWTSecret        string
	JWTAccessExpiry  int
	JWTRefreshExpiry int

	// Security
	RateLimitPerMinute     int
	MaxLoginAttempts       int
	LockoutDurationMinutes int

	// PIN Policy
	PINMinLength int
	PINMaxLength int

	// OTP
	OTPExpiry int
	OTPLength int
}

var AppConfig *Config

func LoadConfig() {
	godotenv.Load()

	AppConfig = &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/cbs_simulator?sslmode=disable"),
		Environment: getEnv("ENVIRONMENT", "development"),

		JWTSecret:        getEnv("JWT_SECRET", "cbs-simulator-secret-key-change-in-production"),
		JWTAccessExpiry:  getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 15),
		JWTRefreshExpiry: getEnvInt("JWT_REFRESH_EXPIRY_HOURS", 168),

		RateLimitPerMinute:     getEnvInt("RATE_LIMIT_PER_MINUTE", 60),
		MaxLoginAttempts:       getEnvInt("MAX_LOGIN_ATTEMPTS", 3),
		LockoutDurationMinutes: getEnvInt("LOCKOUT_DURATION_MINUTES", 30),

		PINMinLength: getEnvInt("PIN_MIN_LENGTH", 6),
		PINMaxLength: getEnvInt("PIN_MAX_LENGTH", 6),

		OTPExpiry: getEnvInt("OTP_EXPIRY_MINUTES", 5),
		OTPLength: getEnvInt("OTP_LENGTH", 6),
	}
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
