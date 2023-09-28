package config

import (
	"fmt"
	"os"

	"github.com/alhazenlabs/procure-hub/backend/internal/logger"
)

var (
	PG_HOST     = GetEnv("PG_HOST", "0.0.0.0")
	PG_PORT     = GetEnv("PG_PORT", "5432")
	PG_USER     = GetEnv("PG_USER", "postgres")
	PG_DB       = GetEnv("PG_DB", "postgres")
	PG_PASSWORD = GetEnv("PG_PASSWORD", "admin@123")
	PG_SSLMODE  = GetEnv("PG_SSLMODE", "disable")
)

func GetEnv(key, defaultValue string) string {
	// Try to retrieve the environment variable
	value := os.Getenv(key)

	// Check if the environment variable is set
	if value == "" {
		logger.Info(fmt.Sprintf("key: %s is not set, using fallback value: %s", key, value))
		value = defaultValue
	}
	return value
}
