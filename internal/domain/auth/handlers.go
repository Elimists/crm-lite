package auth

import (
	"crm-lite/internal/json"
	"html/template"
	"net/http"
)

type handler struct {
	service Service
	views   *template.Template
}

func NewHandler(service Service, views *template.Template) *handler {
	return &handler{
		service: service,
		views:   views,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

	username, password := r.FormValue("username"), r.FormValue("password")

	token, err := h.service.Authenticate(r.Context(), username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"token": token})
}

func (h *handler) GetLoginPage(w http.ResponseWriter, r *http.Request) {

	err := h.views.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), 500)
	}
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {

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
