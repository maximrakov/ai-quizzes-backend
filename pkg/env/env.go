package env

import (
	"os"
)

func GetEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
