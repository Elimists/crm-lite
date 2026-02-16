package main

import (
	"context"
	"crm-lite/db"
	"crm-lite/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Repo   *db.Queries
	Logger *log.Logger
}

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	app := &Application{
		db.New(pool),
		l,
	}
	app.Logger.Printf("Connected to DB. Max connections allowed: %d", pool.Config().MaxConns)
	app.Logger.Println("Application started successfully")

	http.Handle("/", http.HandlerFunc(HomeHandler))
	http.Handle("POST /contacts", utils.ClientMiddleware(http.HandlerFunc(app.HandleCreateContact)))

	app.Logger.Println("Server started on :8080")
	app.Logger.Fatal(http.ListenAndServe(":8080", nil))

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"name":    "crm-lite",
		"version": "2.1",
		"status":  "running",
	})

}

type ContactRequest struct {
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	Message      string      `json:"message"`
	SourceDomain string      `json:"source_domain"`
	Slug         string      `json:"tenant_slug"`
}

func (app *Application) HandleCreateContact(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("Incoming request from origin:", r.Header.Get("Origin"))

	if r.Body == nil {
		http.Error(w, "empty body", 400)
		return
	}
	defer r.Body.Close()

	var payload ContactRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = app.Repo.InsertContact(r.Context(), db.InsertContactParams{
		Name:         payload.Name,
		Email:        payload.Email,
		Phone:        payload.Phone,
		Message:      payload.Message,
		SourceDomain: payload.SourceDomain,
		Slug:         payload.Slug,
	})
	if err != nil {
		app.Logger.Printf("Internal server error %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	/*
		if err := utils.SendEmailNotification(r.Context(), created); err != nil {
			app.Logger.Printf("Failed to send alert email to %s. Error details: %v", created.Email, err)
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
