package controllers

import (
	"crm-lite/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func DisplayLoginPage(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("token"); err == nil {
		if valid, _ := utils.ValidateToken(cookie.Value); valid {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	err := templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "missing required values", http.StatusNotAcceptable)
		return
	}

	u, err := c.QueryUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[INFO:Login] attempt for non-existent email: %s", email)
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}
		log.Printf("[ERROR:Login->QueryUserByEmail] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if u == nil {
		log.Printf("email not found in db : %s", email)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	if err != nil {
		log.Printf("invalid credentials for: %s", email)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	t, err := c.QueryTenant(r.Context(), u.TenantId)
	if err != nil {
		log.Printf("[ERROR:Login->QueryTenant] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tokenString, err := utils.CreateUserToken(
		u.UserName,
		t.Slug,
		u.Roles,
		u.Scopes,
	)
	if err != nil {
		log.Println("[ERR:CreateUserToken] token was nil")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((24*time.Hour - time.Hour).Seconds()),
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {

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
