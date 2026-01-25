package controller

import (
	"crm-lite/db"
	"crm-lite/models"
	"crm-lite/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func CreateContact(dbConn *sql.DB, w http.ResponseWriter, r *http.Request) {

	log.Println("incoming request from origin:", r.Header.Get("Origin"))

	var c models.Contact

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	db.InsertContact(dbConn, c)
	log.Println("new contact created for ", c.Email)

	if err := utils.SendEmailNotification(r.Context(), c); err != nil {
		log.Println("failed to send alert email: ", err)
	}

	w.WriteHeader(http.StatusCreated)
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}
