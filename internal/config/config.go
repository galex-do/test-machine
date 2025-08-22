package config

import (
        "os"
)

// Config holds application configuration
type Config struct {
        DatabaseURL  string
        DatabaseName string
        Port         string
}

// Load loads configuration from environment variables
func Load() *Config {
        return &Config{
                DatabaseURL:  getEnv("DATABASE_URL", "postgres://testmgmt:password@localhost:5432/testmanagement?sslmode=disable"),
                DatabaseName: getEnv("DATABASE_NAME", "testmanagement"),
                Port:         getEnv("PORT", "5000"),
        }
}

func getEnv(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
                return value
        }
        return defaultValue
}