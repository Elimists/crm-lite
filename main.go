package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if err := loadAllowedOrigins("allowed_origins.json"); err != nil {
		log.Fatal("[INIT_ERR] Failed to load allowed origins: ", err)
	}

	required := []string{"MAILTRAP_API_TOKEN", "MAILTRAP_ENDPOINT"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			log.Fatalf("[INIT_ERR] Environment variable %s is not set. Please set it before running the application.", v)
		}
	}

	var err error
	db, err = sql.Open("sqlite3", "./database/contacts.db")
	if err != nil {
		log.Fatal("[INIT_ERR] Failed to open database:", err)
	}

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
	if _, err := db.Exec(sqlStmt); err != nil {
		log.Fatal("[INIT_ERR] Failed to create table:", err)
	}

	log.Println("Initialization complete: DB + Allowed Origins + SMTP OK.")

}

func sendMailtrapEmail(c Contact) error {
	apiToken := os.Getenv("MAILTRAP_API_TOKEN")

	payloadStr := fmt.Sprintf(`{
		"to": [
			{
				"email":"info@weather-wizards.com",
				"name":"James Belanger"
			}
		],
		"bcc": [
			{
				"email": "pandey.pran@gmail.com",
				"name": "Pran Pandey"
			}
		],
		"from": {
			"email":"no-reply@proreact.dev",
			"name":"Weather Wizards Contact Form"
		},
		"subject":"New Contact Form Submission",
		"text":"New contact submitted:\n\nName: %s\nEmail: %s\nPhone: %s\nMessage: %s",
		"category":"Weather Wizards Contact Form"
	}`, c.Name, c.Email, c.Phone, c.Message)

	payload := strings.NewReader(payloadStr)
	req, err := http.NewRequest("POST", os.Getenv("MAILTRAP_ENDPOINT"), payload)
	if err != nil {
		fmt.Println("[MAIL_ERROR] Failed to create request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Token", apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[MAIL_ERROR] Failed to send request:", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[MAIL_ERROR] Failed to read response body:", err)
		return err
	}
	return nil
}

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

	source := r.Header.Get("Origin")
	if source == "" {
		log.Println("incoming request is missing origin header")
		http.Error(w, "missing origin header", http.StatusBadRequest)
		return
	}

	if !allowedOrigins[source] {
		log.Println("origin not allowed:", source)
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}
	log.Println("incoming request from origin:", source)

	var c Contact

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
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
	log.Println("new contact created for ", c.Email)

	if err := sendMailtrapEmail(c); err != nil {
		log.Println("failed to send alert email: ", err)
	}
	log.Println("alert email sent for ", c.Email)

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

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/contact", ContactHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
