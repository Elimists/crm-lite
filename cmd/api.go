package main

import (
	"context"
	"crm-lite/internal/domain/auth"
	"crm-lite/internal/domain/contacts"
	"crm-lite/internal/domain/health"
	"crm-lite/internal/env"
	"crm-lite/internal/jwt"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Minute))

	staticFS, _ := fs.Sub(app.static, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	healthService := health.NewService(
		app.config.App.Name,
		app.config.App.Slug,
		app.config.App.Version)
	healthHandler := health.NewHandler(healthService)
	r.Get("/health", healthHandler.GetHealth)

	contactService := contacts.NewService(app.db)
	contactHandler := contacts.NewHandler(contactService)
	r.Route("/contacts", func(r chi.Router) {
		r.Post("/", contactHandler.CreateContact)
		r.Group(func(r chi.Router) {
			r.Use(app.AuthMiddleware)
			r.Get("/", contactHandler.GetContacts)
			r.Get("/{id}", contactHandler.GetContact)
		})

	})

	authService := auth.NewService(app.db, app.authenticator)
	authHandler := auth.NewHandler(authService, app.views)
	r.Route("/login", func(r chi.Router) {
		r.Get("/", authHandler.GetLoginPage)
		r.Post("/", authHandler.Login)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:    app.config.App.Addr,
		Handler: h,
	}

	log.Printf("server has started at addr %s", app.config.App.Addr)
	return server.ListenAndServe()
}

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := app.authenticator.ValidateToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var claimsKey contextKey = contextKey(app.config.App.Slug + "_claims")
		ctx := context.WithValue(r.Context(), claimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey string
type application struct {
	config        *env.Config
	db            *pgxpool.Pool
	views         *template.Template
	authenticator *jwt.Authenticator
	static        embed.FS
}
