package main

import (
	"crm-lite/internal/contacts"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	contactService := contacts.NewService()
	contactHandler := contacts.NewHandler(contactService)
	r.Route("/contacts", func(r chi.Router) {
		r.Post("/", contactHandler.CreateContact)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:    app.config.addr,
		Handler: h,
	}

	log.Printf("server has started at addr %s", app.config.addr)
	return server.ListenAndServe()
}

type application struct {
	config config
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
