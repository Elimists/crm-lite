package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Warning: The %s environment variable was not set, using default %d", key, fallback)
		return fallback
	}
	converted, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Warning: %s '%s' is not an integer, using default %d", key, val, fallback)
		return fallback
	}
	return converted
}
