package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// CORS
	FrontendURL string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Logging
	LogLevel string
}

// Load returns the application configuration loaded from environment variables
func Load() *Config {

	return &Config{
		Port:          getEnv("PORT", "8080").String(),
		Environment:   getEnv("ENVIRONMENT", "development").String(),
		DatabaseURL:   getEnv("DATABASE_URL", "").String(),
		JWTSecret:     getEnv("JWT_SECRET", "change-this-secret").String(),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h").AsDuration(),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:3000").String(),
		RedisHost:     getEnv("REDIS_HOST", "localhost").String(),
		RedisPort:     getEnv("REDIS_PORT", "6379").String(),
		RedisPassword: getEnv("REDIS_PASSWORD", "").String(),
		LogLevel:      getEnv("LOG_LEVEL", "info").String(),
	}
}

func getEnv(key, defaultValue string) *EnvValues {
	if value := os.Getenv(key); value != "" {
		return &EnvValues{value: value}
	}
	return &EnvValues{value: defaultValue}
}

type EnvValues struct {
	value string
}

func (e *EnvValues) String() string {
	return e.value
}

func (e *EnvValues) AsInt() (int, error) {
	// Parse as int, return error if unable
	var result int
	_, err := fmt.Sscanf(e.value, "%d", &result)
	return result, err
}

func (e *EnvValues) AsBool() (bool, error) {
	// Parse as bool, return error if unable
	return strconv.ParseBool(e.value)
}

func (e *EnvValues) AsDuration() time.Duration {
	// Parse as duration, return error if unable
	val, _ := time.ParseDuration(e.value)
	return val
}
