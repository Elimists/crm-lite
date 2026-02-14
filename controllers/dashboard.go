package controllers

import (
	"crm-lite/utils"
	"net/http"
)

func DisplayDashboardPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	valid, _ := utils.ValidateToken(cookie.Value)

	if !valid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Sat, 01 Jan 2000 00:00:00 GMT")

	err = templates.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
