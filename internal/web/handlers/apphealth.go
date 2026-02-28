package handlers

import (
	"crm-lite/internal/domain/health"
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	service health.Service
}

func NewHealthHandler(service health.Service) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	data := h.service.GetHealth()

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
