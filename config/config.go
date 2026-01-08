// Package config handles the env variables and configuration needed to run the app
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	WeatherAPI    EnvKey = "WEATHER_API"
	RedisAddr     EnvKey = "REDIS_ADDR"
	RedisPassword EnvKey = "REDIS_PASSWORD"
	ServerPort    EnvKey = "SERVER_PORT"
)

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

// Load loads the .env file and validates required configuration
func Load() error {
	if err := godotenv.Load(".env"); err != nil {
		// It's OK if .env doesn't exist - environment variables might be set
		fmt.Printf("Warning: .env file not loaded: %v\n", err)
	}

	// Validate required configuration
	if err := validateRequiredConfig(); err != nil {
		return err
	}

	return nil
}

// validateRequiredConfig checks that all required environment variables are set
func validateRequiredConfig() error {
	requiredKeys := []EnvKey{WeatherAPI}

	for _, key := range requiredKeys {
		if value := key.GetValue(); value == "" {
			return fmt.Errorf("required configuration missing: %s", key)
		}
	}

	return nil
}
