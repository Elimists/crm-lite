package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var allowedOrigins map[string]bool

func loadAllowedOrigins(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var origins []string
	err = json.Unmarshal(data, &origins)
	if err != nil {
		return err
	}

	allowedOrigins = make(map[string]bool)
	for _, o := range origins {
		allowedOrigins[o] = true
	}

	return nil
}

var db *sql.DB

type Contact struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Message      string `json:"message"`
	SourceDomain string `json:"source_domain"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"name":    "crm-lite",
		"version": "1.0.0",
		"status":  "running",
	})

}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateContact(w, r)
		/*
			case http.MethodGet:
				GetContacts(w, r)
		*/
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateContact(w http.ResponseWriter, r *http.Request) {

	var c Contact

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	source := r.Header.Get("Origin")
	log.Println("Incoming Origin:", source)
	if source == "" {
		http.Error(w, "missing origin header", http.StatusBadRequest)
		return
	}

	if !allowedOrigins[source] {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}

	stmt, err := db.Prepare("INSERT INTO contacts(name, email, phone, message, source_domain) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "database error (prep)", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.Email, c.Phone, c.Message, source)
	if err != nil {
		http.Error(w, "batabase error (exec)", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	pageSize := 20
	page := 1

	// Get total count
	var totalCount int
	err := db.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&totalCount)
	if err != nil {
		http.Error(w, "database error (sele)"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate max page number
	maxPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if maxPage == 0 {
		maxPage = 1
	}

	// Parse page from query parameter if totalCount > pageSize
	if totalCount > pageSize {
		if p := r.URL.Query().Get("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 {
				page = n
			}
		}
	}

	// Cap page number to maxPage
	if page > maxPage {
		page = maxPage
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT uuid, name, email, phone, message, source_domain, status, created_at, updated_at
		FROM contacts
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		http.Error(w, "Database error "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(&c.UUID, &c.Name, &c.Email, &c.Phone, &c.Message, &c.SourceDomain, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			http.Error(w, "Database error "+err.Error(), http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, c)
	}

	response := map[string]interface{}{
		"count":      len(contacts),
		"totalCount": totalCount,
		"page":       page,
		"pageSize":   pageSize,
		"contacts":   contacts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

func main() {
	err := loadAllowedOrigins("allowed_origins.json")
	if err != nil {
		log.Fatal("Failed to load allowed origins:", err)
	}

	db, err = sql.Open("sqlite3", "./database/contacts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL DEFAULT (lower(hex(randomblob(8)))),
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		phone TEXT,
		message TEXT NOT NULL,
		source_domain TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'new',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Contacts table initialized.")

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/contact", ContactHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
