package contacts

import (
	"crm-lite/internal/json"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateContact(w http.ResponseWriter, r *http.Request) {

	contact := []string{"hi", "world"}
	json.Write(w, 200, contact)
}
