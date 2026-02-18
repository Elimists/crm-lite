package main

import (
	"context"
	"crm-lite/internal/domain/contacts"
	"crm-lite/internal/jwt"
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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

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

	return r
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
		var claimsKey contextKey = contextKey(app.config.slug + "_claims")
		ctx := context.WithValue(r.Context(), claimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:    app.config.addr,
		Handler: h,
	}

	log.Printf("server has started at addr %s", app.config.addr)
	return server.ListenAndServe()
}

type contextKey string
type application struct {
	config        config
	db            *pgxpool.Pool
	authenticator jwt.Authenticator
}

type config struct {
	name    string
	slug    string
	version string
	addr    string
	db      dbConfig
}

type dbConfig struct {
	host        string
	port        string
	user        string
	pass        string
	name        string
	sslmode     string
	maxConn     int32
	maxIdleTime time.Duration
}
