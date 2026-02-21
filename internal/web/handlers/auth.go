package handlers

import (
	"crm-lite/internal/domain/auth"
	"crm-lite/internal/json"
	"html/template"
	"net/http"
)

type AuthHandler struct {
	service auth.Service
	views   *template.Template
}

func NewAuthHandler(service auth.Service, views *template.Template) *AuthHandler {
	return &AuthHandler{
		service: service,
		views:   views,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	username, password := r.FormValue("username"), r.FormValue("password")

	token, err := h.service.Authenticate(r.Context(), username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) GetLoginPage(w http.ResponseWriter, r *http.Request) {

	err := h.views.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), 500)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
