package controllers

import (
	"crm-lite/db"
	"crm-lite/models"
	"crm-lite/utils"
	"encoding/json"
	"log"
	"net/http"
)

func ContactHandler(db *db.Database, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateContact(db, w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateContact(db *db.Database, w http.ResponseWriter, r *http.Request) {

	log.Println("incoming request from origin:", r.Header.Get("Origin"))

	var contact models.Contact

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	db.InsertContact(r.Context(), contact)
	log.Println("new contact created for ", contact.Email)

	if err := utils.SendEmailNotification(contact); err != nil {
		log.Println("failed to send alert email: ", err)
	}

	w.WriteHeader(http.StatusCreated)
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}
