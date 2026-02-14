package controllers

import (
	"crm-lite/db"
	"crm-lite/utils"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantController struct {
	db *db.Database
}

func NewTenantController(database *db.Database) *TenantController {
	return &TenantController{db: database}
}

func (c *TenantController) GetTenant(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	valid, _ := utils.ValidateToken(token.Value)
	if !valid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	t, err := c.QueryTenant(r.Context(), 1)
	if err != nil {
		http.Error(w, "db query error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("dashboard.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, t)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}

func (c *Controller) GetTenants(w http.ResponseWriter, r *http.Request) {

	token, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	valid, claims := utils.ValidateToken(token.Value)
	if !valid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if claims["tenant"] != "proreact" {
		Logout(w, r)
		return
	}

	t, err := c.QueryTenants(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("tenants.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.Execute(w, t)
	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}

func (c *Controller) CreateTenant(dbConn *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) PatchTenant(dbConn *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) DeleteTenant(dbConn *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

}
