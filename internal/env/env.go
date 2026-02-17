package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	s := os.Getenv(key)
	if s == "" {
		return fallback
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}

	return v
}

func GetDuration(key, fallback string) time.Duration {
	s := os.Getenv(key)
	if s == "" {
		s = fallback
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		// If the env is invalid (like "15minutes" instead of "15m"),
		// fallback to a safe default.
		d, _ = time.ParseDuration(fallback)
	}
	return d
}
