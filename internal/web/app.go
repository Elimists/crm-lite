package web

import (
	"crm-lite/internal/env"
	"crm-lite/internal/jwt"
	"embed"
	"html/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config        *env.Config
	DB            *pgxpool.Pool
	Authenticator *jwt.Authenticator // Using the pointer as discussed!
	Views         *template.Template
	Static        embed.FS
}

func NewApplication(cfg *env.Config, db *pgxpool.Pool, auth *jwt.Authenticator, views *template.Template, static embed.FS) *Application {
	return &Application{
		Config:        cfg,
		DB:            db,
		Authenticator: auth,
		Views:         views,
		Static:        static,
	}
}
