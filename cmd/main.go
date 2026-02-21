package main

import (
	"context"
	"crm-lite/internal/env"
	"crm-lite/internal/jwt"
	"crm-lite/internal/web"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	// Configuration
	cfg := env.LoadConfig()

	// Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("unable to connect to db", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to database", "db", cfg.DB.Name)

	authenticator := jwt.NewAuthenticator(
		env.GetString("JWT_SECRET", "devdefaultrandomverylongstring1a2b3c4d5e"),
		cfg.App.Slug)

	// Load HTML templates
	views := template.Must(template.ParseFS(
		web.ViewFiles, "views/*/*.html"))

	/*
		for _, t := range views.Templates() {
			fmt.Println("Available template:", t.Name())
		}

	*/

	// Application
	api := &web.Application{
		Config:        cfg,
		DB:            db,
		Views:         views,
		Authenticator: authenticator,
		Static:        web.StaticFiles,
	}
	if err := api.Run(api.Mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
