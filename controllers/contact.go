package controllers

import (
	"crm-lite/db"
	"crm-lite/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type ContactHandler struct {
	Queries *db.Queries
}

type ContactRequest struct {
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	Phone        pgtype.Text `json:"phone"`
	Message      string      `json:"message"`
	SourceDomain string      `json:"source_domain"`
	Slug         string      `json:"tenant_slug"`
}

func (h *ContactHandler) HandleCreateContact(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming request from origin:", r.Header.Get("Origin"))

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

	c, err := h.Queries.InsertContact(r.Context(), db.InsertContactParams{
		Name:         payload.Name,
		Email:        payload.Email,
		Phone:        payload.Phone,
		Message:      payload.Message,
		SourceDomain: payload.SourceDomain,
		Slug:         payload.Slug,
	})
	if err != nil {
		log.Printf("Internal server error %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t, _ := h.Queries.GetTenantBySlug(r.Context(), payload.Slug)

	go func(t *db.Tenant, c *db.Contact) {
		if err := utils.NotifyTenant(t, c); err != nil {
			log.Printf("Failed to send alert email to tenant %s (%s). Error details: %v", t.Name, t.Email, err)
		}
	}(&t, &c)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
