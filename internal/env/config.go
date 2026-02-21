package env

import (
	"time"
)

type Config struct {
	App struct {
		Name    string
		Slug    string
		Version string
		Addr    string
	}
	DB struct {
		Host        string
		Port        string
		User        string
		Pass        string
		Name        string
		SSLMode     string
		MaxConn     int32
		MaxIdleTime time.Duration
	}
}

func LoadConfig() *Config {
	var cfg Config

	// App Settings
	cfg.App.Name = GetString("APP_NAME", "CRM Lite (DEV)")
	cfg.App.Slug = GetString("APP_SLUG", "crmlite-dev")
	cfg.App.Version = GetString("APP_VERSION", "1.0.0-dev")
	cfg.App.Addr = ":" + GetString("APP_PORT", "8080")

	// DB Settings
	cfg.DB.Host = GetString("DB_HOST", "localhost")
	cfg.DB.Port = GetString("DB_PORT", "5432")
	cfg.DB.User = GetString("DB_USER", "postgres")
	cfg.DB.Pass = GetString("DB_PASSWORD", "postgres")
	cfg.DB.Name = GetString("DB_NAME", cfg.App.Slug) // Defaults to app slug
	cfg.DB.SSLMode = GetString("DB_SSLMODE", "disable")
	cfg.DB.MaxConn = int32(GetInt("DB_MAX_CONN", 25))
	cfg.DB.MaxIdleTime = GetDuration("DB_MAX_IDLE_TIME", "15m")

	return &cfg
}
