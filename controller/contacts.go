package controller

import (
	"crm-lite/db"
	"crm-lite/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func CreateContact(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {

	source := r.Header.Get("Origin")
	if source == "" {
		log.Println("incoming request is missing origin header")
		http.Error(w, "missing origin header", http.StatusBadRequest)
		return
	}

	/* TODO move logic to utils pacakge
	if !allowedOrigins[source] {
		log.Println("origin not allowed:", source)
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}
	*/
	log.Println("incoming request from origin:", source)

	var c models.Contact

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	db.InsertContact(dbConn, c)
	log.Println("new contact created for ", c.Email)

	/* TODO move to utils package
	if err := sendMailtrapEmail(c); err != nil {
		log.Println("failed to send alert email: ", err)
	}
	log.Println("alert email sent for ", c.Email)
	*/

	w.WriteHeader(http.StatusCreated)
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}
