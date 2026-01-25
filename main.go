package main

import (
	"crm-lite/controller"
	"crm-lite/db"
	"crm-lite/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	utils.LoadClientConfigs("clients.json")
	utils.VerifyRequiredEnv()

	var err error
	dbConn, err = db.OpenDB("contacts.db")
	if err != nil {
		log.Fatal("[INIT_ERR] Failed to open database:", err)
	}

	if err := db.SetupTables(dbConn); err != nil {
		log.Fatal("[MIGRATION_ERR] Failed to create tables:", err)
	}

	log.Println("Initialization complete: DB + Allowed Origins + SMTP OK.")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"name":    "crm-lite",
		"version": "2.0",
		"status":  "running",
	})

}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		controller.CreateContact(dbConn, w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	http.Handle("/", http.HandlerFunc(HomeHandler))
	http.Handle("/contact", utils.ClientMiddleware(http.HandlerFunc(ContactHandler)))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
