package contacts

import (
	"crm-lite/internal/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetContact(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	contact, err := h.service.GetContact(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.Write(w, http.StatusOK, contact)
}

func (h *handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("all contacts - not allowed"))
}

func (h *handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var tempContact createContactParams

	if err := json.Read(r, &tempContact); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdContact, err := h.service.CreateContact(r.Context(), tempContact)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdContact)
}
