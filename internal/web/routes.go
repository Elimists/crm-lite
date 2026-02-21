package web

import (
	"crm-lite/internal/domain/auth"
	"crm-lite/internal/domain/contacts"
	"crm-lite/internal/domain/health"
	"crm-lite/internal/web/handlers"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	// Standard Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 60))

	staticFS, err := fs.Sub(app.Static, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// 1. Health
	healthService := health.NewService(
		app.Config.App.Name,
		app.Config.App.Slug,
		app.Config.App.Version)
	healthHandler := handlers.NewHealthHandler(healthService)
	r.Get("/health", healthHandler.GetHealth)

	// 2. Auth
	authService := auth.NewService(app.DB, app.Authenticator)
	authHandler := handlers.NewAuthHandler(authService, app.Views)
	r.Route("/login", func(r chi.Router) {
		r.Get("/", authHandler.GetLoginPage)
		r.Post("/", authHandler.Login)
	})

	// 3. Contacts
	contactService := contacts.NewService(app.DB)
	contactHandler := handlers.NewContactHandler(contactService)
	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware) // defined in middleware.go
		r.Route("/contacts", func(r chi.Router) {
			r.Get("/", contactHandler.GetContacts)
			r.Post("/", contactHandler.CreateContact)
		})
	})

	return r
}

func (app *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.Config.App.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server starting on %s (%s)", app.Config.App.Addr, app.Config.App.Name)

	// This will block until the server stops or errors out
	return srv.ListenAndServe()
}
