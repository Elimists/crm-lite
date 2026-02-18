package auth

import (
	"crm-lite/internal/json"
	"crm-lite/internal/jwt"
	"log"
	"net/http"
)

type Handler struct {
	authenticator *jwt.Authenticator
	service       Service
}

func NewHandler(auth *jwt.Authenticator) *Handler {
	return &Handler{authenticator: auth}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	user, err := h.service.Authenticate(r.Context(), "", "")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.authenticator.CreateToken(
		user.UserName,
		user.TenantSlug,
		user.Roles,
		user.Scopes,
	)
	if err != nil {
		http.Error(w, "error", 500)
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"token": token})
}
