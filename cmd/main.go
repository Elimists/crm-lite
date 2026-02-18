package main

import (
	"context"
	"crm-lite/internal/env"
	"crm-lite/internal/jwt"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	// Configuration
	cfg := &config{
		name:    env.GetString("APP_NAME", "CRM Lite"),
		slug:    env.GetString("APP_SLUG", "crmlite"),
		version: env.GetString("APP_VERSION", "1.0.0-dev"),
		addr:    ":" + env.GetString("APP_PORT", "8080"),
	}
	cfg.db = dbConfig{
		host:        env.GetString("DB_HOST", "localhost"),
		port:        env.GetString("DB_PORT", "5432"),
		user:        env.GetString("DB_USER", "pran"),
		pass:        env.GetString("DB_PASSWORD", "abc@123"),
		name:        env.GetString("DB_NAME", cfg.slug),
		sslmode:     env.GetString("DB_SSLMODE", "disable"),
		maxConn:     int32(env.GetInt("DB_MAX_CONN", 25)),
		maxIdleTime: env.GetDuration("DB_MAX_IDLE_TIME", "15m"),
	}

	// Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.db.user, cfg.db.pass, cfg.db.host, cfg.db.port, cfg.db.name, cfg.db.sslmode)
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("unable to connect to db", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to database", "db", cfg.db.name)

	authenticator := jwt.NewAuthenticator(
		env.GetString("JWT_SECRET", "defaultrandomverylongstring1a2b3c4d5e"),
		cfg.slug)

	// Application
	api := &application{
		config:        *cfg,
		db:            db,
		authenticator: *authenticator,
	}
	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
